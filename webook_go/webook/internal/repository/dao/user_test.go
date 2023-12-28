package dao

import (
	"context"
	"testing"
)

func TestGORMUserDAO_Insert(t *testing.T) {
	type fields struct {
		db *gorm.DB
	}
	type args struct {
		ctx context.Context
		u   User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dao := &GORMUserDAO{
				db: tt.fields.db,
			}
			if err := dao.Insert(tt.args.ctx, tt.args.u); (err != nil) != tt.wantErr {
				t.Errorf("Insert() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
