package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAbs(t *testing.T) {
	tests := []struct {
		name   string
		values float64
		want   float64
	}{
		{
			name:   "-3 to 3",
			values: -3,
			want:   3,
		},
		{
			name:   "-2.000001 to 2.000001",
			values: -2.000001,
			want:   2.000001,
		},
		{
			name:   "-0.000000003 to 0.000000003",
			values: -0.000000003,
			want:   0.000000003,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, Abs(test.values))
			//if res := Abs(test.values); res != test.want {
			//	t.Errorf("Abs() = %f, want %f", res, test.want)
			//}
		})
	}
}

func TestUser_FullName(t *testing.T) {
	tests := []struct {
		name   string
		values User
		want   string
	}{
		{
			name: "1 test",
			values: User{
				FirstName: "Nurlan",
				LastName:  "Ibragimov",
			},
			want: "Nurlan Ibragimov",
		},
		{
			name: "2 test",
			values: User{
				FirstName: "Petr",
				LastName:  "Petrov",
			},
			want: "Petr Petrov",
		},
	}

	for _, test := range tests { // цикл по всем тестам
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, test.values.FullName())
			//if res := test.values.FullName(); res != test.want {
			//	t.Errorf("user.FullName() = %s, want %s", res, test.want)
			//}
		})
	}
}

func TestFamily_AddNew(t *testing.T) {
	f := Family{}

	tests := []struct {
		name         string
		relationship Relationship
		values       Person
		want         bool
	}{
		{
			name:         "add father",
			relationship: Father,
			values: Person{
				FirstName: "Nurlan",
				LastName:  "Ibragimov",
				Age:       31,
			},
			want: false,
		},
		{
			name:         "double add father",
			relationship: Father,
			values: Person{
				FirstName: "Petr",
				LastName:  "Petrov",
				Age:       31,
			},
			want: true,
		},
	}

	for _, test := range tests { // цикл по всем тестам
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, f.AddNew(test.relationship, test.values) != nil)
			//if res := f.AddNew(test.relationship, test.values); (res != nil) != test.want {
			//	t.Errorf("f.AddNew() = %s, want %t", res, test.want)
			//}
		})
	}
}
