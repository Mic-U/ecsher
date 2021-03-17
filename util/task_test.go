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
	}

	for _, c := range cases {
		result := FilterTasksByNames(c.tasks, c.names)
		if len(result) != 1 {
			t.Fatal("len(result) should be 1")
		}
		resultName := ArnToName(result[0])
		if resultName != c.names[0] {
			t.Errorf("expect %v, got %v", c.names[0], resultName)
		}
	}
}
