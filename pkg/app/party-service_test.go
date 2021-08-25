package app

import (
	"fmt"
	"testing"

	mock_party "github.com/areknoster/table-driven-tests-gomock/mocks/pkg/party"
	"github.com/areknoster/table-driven-tests-gomock/pkg/party"
	"github.com/golang/mock/gomock"
)

func TestPartyService_GreetVisitors_NotNiceReturnsError(t *testing.T) {
	// initialize gomock controller
	ctrl := gomock.NewController(t)
	// if not all expectations set on the controller are fulfilled at the end of function, the test will fail!
	defer ctrl.Finish()
	// init structure which implements party.VisitorLister interface
	mockedVisitorLister := mock_party.NewMockVisitorLister(ctrl)
	// mockedVisitorLister called once with party.NiceVisitor argument would return []string{"Peter"}, nil
	mockedVisitorLister.EXPECT().ListVisitors(party.NiceVisitor).Return([]party.Visitor{{"Peter", "TheSmart"}}, nil)
	// mockedVisitorLister called once with party.NotNiceVisitor argument would return nil and error
	mockedVisitorLister.EXPECT().ListVisitors(party.NotNiceVisitor).Return(nil, fmt.Errorf("dummyErr"))
	// mockedVisitorLister implements party.VisitorLister interface, so it can be assigned in PartyService
	sp := &PartyService{
		visitorLister: mockedVisitorLister,
	}
	gotErr := sp.GreetVisitors(false)
	if gotErr == nil {
		t.Errorf("did not get an error")
	}
}

func TestPartyService_GreetVisitors(t *testing.T) {
	type fields struct {
		visitorLister *mock_party.MockVisitorLister
		greeter       *mock_party.MockGreeter
	}
	type args struct {
		justNice bool
	}
	tests := []struct {
		name    string
		prepare func(f *fields)
		args    args
		wantErr bool
	}{
		{
			name: "visitorLister.ListVisitors(party.NiceVisitor) returns error, error expected",
			prepare: func(f *fields) {
				f.visitorLister.EXPECT().ListVisitors(party.NiceVisitor).Return(nil, fmt.Errorf("dummyErr"))
			},
			args:    args{justNice: true},
			wantErr: true,
		},
		{
			name: "visitorLister.ListVisitors(party.NotNiceVisitor) returns error, error expected",
			prepare: func(f *fields) {
				// if given calls do not happen in expected order, the test would fail!
				gomock.InOrder(
					f.visitorLister.EXPECT().ListVisitors(party.NiceVisitor).Return([]party.Visitor{party.Visitor{
						Name:    "Madhan",
						Surname: "Kumar",
					},
					}, nil),
					f.visitorLister.EXPECT().ListVisitors(party.NotNiceVisitor).Return(nil, fmt.Errorf("dummyErr")),
				)
			},
			args:    args{justNice: false},
			wantErr: true,
		},
		{
			name: " name of nice person, 1 name of not-nice person. greeter should be called with a nice person first, then with not-nice person as an argument",
			prepare: func(f *fields) {
				nice := []party.Visitor{
					party.Visitor{
						Name:    "Peter",
						Surname: "Parker",
					},
				}
				notNice := []party.Visitor{
					party.Visitor{
						Name:    "Helo",
						Surname: "Parker",
					},
				}
				gomock.InOrder(
					f.visitorLister.EXPECT().ListVisitors(party.NiceVisitor).Return(nice, nil),
					f.visitorLister.EXPECT().ListVisitors(party.NotNiceVisitor).Return(notNice, nil),
					f.greeter.EXPECT().Hello(nice[0].Name+" "+nice[0].Surname),
					f.greeter.EXPECT().Hello(notNice[0].Name+" "+notNice[0].Surname),
				)
			},
			args:    args{justNice: false},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := fields{
				visitorLister: mock_party.NewMockVisitorLister(ctrl),
				greeter:       mock_party.NewMockGreeter(ctrl),
			}
			if tt.prepare != nil {
				tt.prepare(&f)
			}

			s := &PartyService{
				visitorLister: f.visitorLister,
				greeter:       f.greeter,
			}
			if err := s.GreetVisitors(tt.args.justNice); (err != nil) != tt.wantErr {
				t.Errorf("GreetVisitors() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
