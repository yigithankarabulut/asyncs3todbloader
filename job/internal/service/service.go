package service

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/yigithankarabulut/asyncs3todbloader/job/internal/customerror"
	"github.com/yigithankarabulut/asyncs3todbloader/job/model"
	"github.com/yigithankarabulut/asyncs3todbloader/job/pkg/constant"
	"golang.org/x/sync/errgroup"
	"sync"
)

// CheckIfBucketExists method checks if the bucket exists. If the bucket does not exist, it returns an error.
func (s *service) CheckIfBucketExists(ctx context.Context) error {
	_, err := s.s3Client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: &s.s3Data.BucketName,
	})
	if err != nil {
		return customerror.New(constant.ErrBucketNotFound, true).
			Wrap(fmt.Errorf("service.CheckIfBucketExists: %v", err)).
			AddData(fmt.Sprintf("bucketname: %s", s.s3Data.BucketName))
	}
	return nil
}

// CheckIfObjectExists method checks if the object exists. If the object does not exist, it returns an error.
func (s *service) CheckIfObjectExists(ctx context.Context) error {
	_, err := s.s3Client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: &s.s3Data.BucketName,
		Key:    &s.s3Data.ObjectKey,
	})
	if err != nil {
		return customerror.New(constant.ErrObjectNotFound, true).
			Wrap(fmt.Errorf("service.CheckIfObjectExists: %v", err)).
			AddData(fmt.Sprintf("bucketname: %s objectkey: %s", s.s3Data.BucketName, s.s3Data.ObjectKey))
	}
	return nil
}

// CheckObjectDuplicateAndCreate method checks if the object is duplicate in the database. If the object is duplicate, it returns an error.
// It's looking for ContentType, ContentLength, ETag fields of the object. ETag is the MD5 hash of the object and must be unique.
func (s *service) CheckObjectDuplicateAndCreate(ctx context.Context, out *s3.GetObjectOutput) error {
	if out.ContentType == nil {
		return customerror.New(constant.ErrNilObjectFields, true).
			AddData(fmt.Sprintf("bucketname: %s objectkey: %s out.ContentType is nil", s.s3Data.BucketName, s.s3Data.ObjectKey))
	}
	if out.ContentLength == nil {
		return customerror.New(constant.ErrNilObjectFields, true).
			AddData(fmt.Sprintf("bucketname: %s objectkey: %s out.ContentLength is nil", s.s3Data.BucketName, s.s3Data.ObjectKey))
	}
	if out.ETag == nil {
		return customerror.New(constant.ErrNilObjectFields, true).
			AddData(fmt.Sprintf("bucketname: %s objectkey: %s out.ETag is nil", s.s3Data.BucketName, s.s3Data.ObjectKey))
	}
	var objectDetails model.ObjectInfo
	objectDetails.ContentType = *out.ContentType
	objectDetails.ContentLength = *out.ContentLength
	objectDetails.ETag = *out.ETag
	if err := s.objectInfoStorage.Create(ctx, objectDetails); err != nil {
		return customerror.New(constant.ErrCreateObjectInfo, true).
			Wrap(fmt.Errorf("service.CheckObjectDuplicateAndCreate: %v", err)).
			AddData(fmt.Sprintf("bucketname: %s objectkey: %s", s.s3Data.BucketName, s.s3Data.ObjectKey))
	}
	return nil
}

// GetObjectFromS3 method checks if the bucket exists, if the object exists, gets the object from S3 and checks if the object is duplicate.
// After that, it sends the object to the s3OutChan channel to be read. If an error occurs, it returns the error.
func (s *service) GetObjectFromS3(ctx context.Context) error {
	defer close(s.s3OutChan)
	if err := s.CheckIfBucketExists(ctx); err != nil {
		return err
	}
	if err := s.CheckIfObjectExists(ctx); err != nil {
		return err
	}
	out, err := s.s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &s.s3Data.BucketName,
		Key:    &s.s3Data.ObjectKey,
	})
	if err != nil {
		return customerror.New(constant.ErrGetObjectFailed, true).
			Wrap(fmt.Errorf("service.GetObjectFromS3: %v", err)).
			AddData(fmt.Sprintf("bucketname: %s objectkey: %s err: %s", s.s3Data.BucketName, s.s3Data.ObjectKey, err))
	}
	if err := s.CheckObjectDuplicateAndCreate(ctx, out); err != nil {
		return err
	}
	s.s3OutChan <- out
	return nil
}

// ReadDataFromS3Object method reads the object from the s3OutChan channel and sends the lines to the lineChan channel.
// If an error occurs, closes the lineChan channel and returns the error.
func (s *service) ReadDataFromS3Object(ctx context.Context) error {
	out, ok := <-s.s3OutChan
	defer func() {
		close(s.lineChan)
		if out != nil {
			if err := out.Body.Close(); err != nil {
				s.logger.Error(err.Error())
			}
		}
	}()
	if !ok {
		return customerror.New(constant.ErrChannelClosed, true).
			Wrap(fmt.Errorf("service.ReadDataFromS3Object: %v", constant.ErrChannelClosed)).
			AddData("s3OutChan is closed")
	}
	s.logger.Info(fmt.Sprintf("Start reading data from %s", s.s3Data.ObjectKey))
	scanner := bufio.NewScanner(out.Body)
	for scanner.Scan() {
		s.lineChan <- scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		return customerror.New(constant.ErrFileScanFailed, true).
			Wrap(fmt.Errorf("service.ReadDataFromS3Object: %v", err)).
			AddData(fmt.Sprintf("bucketname: %s objectkey: %s line: %s", s.s3Data.BucketName, s.s3Data.ObjectKey, scanner.Text()))
	}
	return nil
}

// HandleLines method reads the lines from the lineChan channel and converts them to the product model.
// Service has a lineHandlerWorkerCount field that determines how many goroutines will be created to handle the lines.
// After that, it sends the product to the productChan channel to be written to the database.
// If an error occurs, closes the productChan channel and returns the error.
func (s *service) HandleLines(ctx context.Context) error {
	wg := sync.WaitGroup{}
	defer close(s.productChan)

	startLine, ok := <-s.lineChan
	if !ok {
		return customerror.New(constant.ErrChannelClosed, true).
			Wrap(fmt.Errorf("service.HandleLines: %v", constant.ErrChannelClosed)).
			AddData("lineChan is closed")
	}
	var product model.Product
	if err := json.Unmarshal([]byte(startLine), &product); err != nil {
		s.logger.Error(fmt.Sprintf("service.HandleLines unmarshal err: %v", err))
	}
	s.productChan <- product
	for i := 0; i < s.lineHandlerWorkerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for line := range s.lineChan {
				var product model.Product
				if err := json.Unmarshal([]byte(line), &product); err != nil {
					s.logger.Error(fmt.Sprintf("service.HandleLines unmarshal err: %v", err))
					continue
				}
				s.productChan <- product
			}
		}()
	}
	wg.Wait()
	return nil
}

// WriteDataToDb method reads the products from the productChan channel and writes them to the database.
// If the channel is closed, to avoid running workers unnecessarily and to log this situation.
// Naturally, we process the first data manually because if there is only 1 data, the channel is closed.
// Same situation have to be handled in HandleLines method.
// Service has a dbWriteWorkerCount field that determines how many goroutines will be created to write the products to the database.
// If an error occurs, returns the error. If the product not written to the database, logs the error.
func (s *service) WriteDataToDb(ctx context.Context) error {
	startProduct, ok := <-s.productChan
	if !ok {
		return customerror.New(constant.ErrChannelClosed, true).
			Wrap(fmt.Errorf("service.WriteDataToDb: %v", constant.ErrChannelClosed)).
			AddData("productChan is closed")
	}
	s.logger.Info(fmt.Sprintf("Start writing data to db"))
	if err := s.productStorage.Create(ctx, startProduct); err != nil {
		var ce *customerror.Error
		if errors.As(err, &ce) {
			message := ce.Message
			if ce.Data != nil {
				data, ok := ce.Data.(string)
				if ok {
					message += ", " + data
				}
				if ce.Loggable {
					s.logger.Error(message)
				}
			}
		}
	}
	wg := sync.WaitGroup{}
	for i := 0; i < s.dbWriteWorkerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for product := range s.productChan {
				if err := s.productStorage.Create(ctx, product); err != nil {
					var ce *customerror.Error
					if errors.As(err, &ce) {
						message := ce.Message
						if ce.Data != nil {
							data, ok := ce.Data.(string)
							if ok {
								message += ", " + data
							}
							if ce.Loggable {
								s.logger.Error(message)
							}
						}
					}
				}
			}
		}()
	}
	wg.Wait()
	return nil
}

// For each S3 object to be read, a goroutine comes to the Run method and runs the methods in funcArr concurrently.
func (s *service) Run() error {
	s.logger.Info(fmt.Sprintf("Start processing %s", s.s3Data.ObjectKey))
	funcArr := []func(ctx context.Context) error{
		s.GetObjectFromS3,
		s.ReadDataFromS3Object,
		s.HandleLines,
		s.WriteDataToDb,
	}
	g, ctx := errgroup.WithContext(context.Background())
	for _, f := range funcArr {
		f := f
		g.Go(func() error {
			return f(ctx)
		})
	}
	if err := g.Wait(); err != nil {
		return err
	}
	return nil
}
