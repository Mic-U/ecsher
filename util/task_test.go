package util

import "testing"

func TestFilterTasksByNames(t *testing.T) {
	cases := []struct {
		tasks []string
		names []string
	}{
		{
			tasks: []string{
				"arn:aws:ecs:region:account:task/cluster-name/abcdefg",
				"arn:aws:ecs:region:account:task/cluster-name/1234567",
			},
			names: []string{
				"abcdefg",
			},
		},
		{
			tasks: []string{
				"arn:aws:ecs:region:account:task/cluster-name/abcdefg",
				"arn:aws:ecs:region:account:task/cluster-name/1234567",
			},
			names: []string{
				"1234567",
			},
		},
		{
			tasks: []string{
				"arn:aws:ecs:region:account:task/cluster-name/abcdefg",
				"arn:aws:ecs:region:account:task/cluster-name/1234567",
			},
			names: []string{},
		},
	}

	for _, c := range cases {
		result := FilterTasksByNames(c.tasks, c.names)
		if len(c.names) == 0 {
			if len(result) != 2 {
				t.Fatal("len(result) should be 2")
			}
			continue
		}
		if len(result) != len(c.names) {
			t.Fatal("len(result) should be 1")
		}
		resultName := ArnToName(result[0])
		if resultName != c.names[0] {
			t.Errorf("expect %v, got %v", c.names[0], resultName)
		}
	}
}
