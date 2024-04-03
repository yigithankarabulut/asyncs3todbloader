package productservice_test

import (
	"context"
	dto "github.com/yigithankarabulut/asyncs3todbloader/microservice/internal/dto/product"
	productrepository "github.com/yigithankarabulut/asyncs3todbloader/microservice/internal/repository/product"
	productservice "github.com/yigithankarabulut/asyncs3todbloader/microservice/internal/service/product"
	"github.com/yigithankarabulut/asyncs3todbloader/microservice/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"reflect"
	"testing"
)

func Test_service_GetProduct(t *testing.T) {
	type fields struct {
		repository productrepository.ProductRepository
	}
	type args struct {
		ctx context.Context
		req dto.GetProductRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    dto.ProductResponse
		wantErr bool
	}{
		{
			name: "Case repository returns error",
			fields: fields{
				repository: &mockProductRepository{
					getErr: ErrGetProductByID,
				},
			},
			args: args{
				ctx: context.Background(),
				req: dto.GetProductRequest{
					ID: 1,
				},
			},
			want:    dto.ProductResponse{},
			wantErr: true,
		},
		{
			name: "Case repository returns product successfully",
			fields: fields{
				repository: &mockProductRepository{
					getErr: nil,
					getRes: model.Product{
						UID:         primitive.NewObjectID(),
						ID:          1,
						Title:       "test title",
						Price:       100,
						Description: "test description",
						Category:    "test category",
						Brand:       "test brand",
						Url:         "test url",
					},
				},
			},
			args: args{
				ctx: context.Background(),
				req: dto.GetProductRequest{
					ID: 1,
				},
			},
			want: dto.ProductResponse{
				ID:          1,
				Title:       "test title",
				Price:       100,
				Description: "test description",
				Category:    "test category",
				Brand:       "test brand",
				Url:         "test url",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := productservice.New(productservice.WithRepository(tt.fields.repository))
			got, err := s.GetProduct(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetProduct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetProduct() got = %v, want %v", got, tt.want)
			}
		})
	}
}
