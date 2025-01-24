package uuid_test

import (
	"testing"

	vectorstore_uuid "github.com/aria3ppp/rag-server/internal/vectorstore/infras/uuid"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

func Test_UuidIDGenerator_NewID(t *testing.T) {
	t.Parallel()

	type want struct {
		err bool
	}

	type testCase struct {
		name string
		want want
	}
	testCases := []testCase{
		{
			name: "ok",
			want: want{
				err: false,
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			idGenerator := vectorstore_uuid.NewIDGenerator()

			uuidString, err := idGenerator.NewID()
			if (err != nil) != tt.want.err {
				t.Fatal(cmp.Diff(err, nil))
			}

			_, err = uuid.Parse(uuidString)
			if err != nil {
				t.Fatal(cmp.Diff(err, nil))
			}
		})
	}
}
