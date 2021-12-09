package main

type testCase struct {
	testDescription string
	requiredClaims  map[string]interface{}
	tokenClaims     map[string]interface{}
	expectedResult  bool
}

var testCases = []testCase{
	{
		testDescription: "both are nil",
		requiredClaims:  nil,
		tokenClaims:     nil,
		expectedResult:  true,
	},
	{
		testDescription: "both are empty",
		requiredClaims:  map[string]interface{}{},
		tokenClaims:     map[string]interface{}{},
		expectedResult:  true,
	},
	{
		testDescription: "required claims are nil",
		requiredClaims:  nil,
		tokenClaims: map[string]interface{}{
			"foo": "bar",
		},
		expectedResult: true,
	},
	{
		testDescription: "required claims are empty",
		requiredClaims:  map[string]interface{}{},
		tokenClaims: map[string]interface{}{
			"foo": "bar",
		},
		expectedResult: true,
	},
	{
		testDescription: "token claims are nil",
		requiredClaims: map[string]interface{}{
			"foo": "bar",
		},
		tokenClaims:    nil,
		expectedResult: false,
	},
	{
		testDescription: "token claims are empty",
		requiredClaims: map[string]interface{}{
			"foo": "bar",
		},
		tokenClaims:    map[string]interface{}{},
		expectedResult: false,
	},
	{
		testDescription: "required is string, token is int",
		requiredClaims: map[string]interface{}{
			"foo": "bar",
		},
		tokenClaims: map[string]interface{}{
			"foo": 1337,
		},
		expectedResult: false,
	},
	{
		testDescription: "matching with string",
		requiredClaims: map[string]interface{}{
			"foo": "bar",
		},
		tokenClaims: map[string]interface{}{
			"foo": "bar",
		},
		expectedResult: true,
	},
	{
		testDescription: "matching with string and int",
		requiredClaims: map[string]interface{}{
			"foo": "bar",
			"bar": 1337,
		},
		tokenClaims: map[string]interface{}{
			"foo": "bar",
			"bar": 1337,
		},
		expectedResult: true,
	},
	{
		testDescription: "matching with string and int in different orders",
		requiredClaims: map[string]interface{}{
			"foo": "bar",
			"bar": 1337,
		},
		tokenClaims: map[string]interface{}{
			"bar": 1337,
			"foo": "bar",
		},
		expectedResult: true,
	},
	{
		testDescription: "matching with string, int and float",
		requiredClaims: map[string]interface{}{
			"foo": "bar",
			"bar": 1337,
			"baz": 13.37,
		},
		tokenClaims: map[string]interface{}{
			"foo": "bar",
			"bar": 1337,
			"baz": 13.37,
		},
		expectedResult: true,
	},
	{
		testDescription: "not matching with string, int and float",
		requiredClaims: map[string]interface{}{
			"foo": "bar",
			"bar": 1337,
			"baz": 13.37,
		},
		tokenClaims: map[string]interface{}{
			"foo": "bar",
			"bar": 1337,
			"baz": 12.27,
		},
		expectedResult: false,
	},
	{
		testDescription: "matching slice",
		requiredClaims: map[string]interface{}{
			"foo": "bar",
			"bar": 1337,
			"baz": []string{"foo"},
		},
		tokenClaims: map[string]interface{}{
			"foo": "bar",
			"bar": 1337,
			"baz": []string{"foo"},
		},
		expectedResult: true,
	},
	{
		testDescription: "matching slice with multiple values",
		requiredClaims: map[string]interface{}{
			"oof": []string{"foo", "bar"},
		},
		tokenClaims: map[string]interface{}{
			"oof": []string{"foo", "bar", "baz"},
		},
		expectedResult: true,
	},
	{
		testDescription: "required slice contains in token slice",
		requiredClaims: map[string]interface{}{
			"foo": "bar",
			"bar": 1337,
			"baz": []string{"foo"},
		},
		tokenClaims: map[string]interface{}{
			"foo": "bar",
			"bar": 1337,
			"baz": []string{"foo", "bar", "baz"},
		},
		expectedResult: true,
	},
	{
		testDescription: "not matching slice",
		requiredClaims: map[string]interface{}{
			"foo": "bar",
			"bar": 1337,
			"baz": []string{"foo"},
		},
		tokenClaims: map[string]interface{}{
			"foo": "bar",
			"bar": 1337,
			"baz": []string{"bar"},
		},
		expectedResult: false,
	},
	{
		testDescription: "matching map",
		requiredClaims: map[string]interface{}{
			"foo": map[string]string{
				"foo": "bar",
			},
		},
		tokenClaims: map[string]interface{}{
			"foo": map[string]string{
				"foo": "bar",
			},
		},
		expectedResult: true,
	},
	{
		testDescription: "matching map with multiple values",
		requiredClaims: map[string]interface{}{
			"foo": map[string]string{
				"foo": "bar",
				"bar": "foo",
			},
		},
		tokenClaims: map[string]interface{}{
			"foo": map[string]string{
				"a":   "b",
				"foo": "bar",
				"bar": "foo",
				"c":   "d",
			},
		},
		expectedResult: true,
	},
	{
		testDescription: "matching map with multiple keys in token claims",
		requiredClaims: map[string]interface{}{
			"foo": map[string]string{
				"foo": "bar",
			},
		},
		tokenClaims: map[string]interface{}{
			"foo": map[string]string{
				"a":   "b",
				"foo": "bar",
				"c":   "d",
			},
		},
		expectedResult: true,
	},
	{
		testDescription: "not matching map",
		requiredClaims: map[string]interface{}{
			"foo": map[string]string{
				"foo": "bar",
			},
		},
		tokenClaims: map[string]interface{}{
			"foo": map[string]int{
				"foo": 1337,
			},
		},
		expectedResult: false,
	},
	{
		testDescription: "matching map with string slice",
		requiredClaims: map[string]interface{}{
			"foo": map[string][]string{
				"foo": {"bar"},
			},
		},
		tokenClaims: map[string]interface{}{
			"foo": map[string][]string{
				"foo": {"foo", "bar", "baz"},
			},
		},
		expectedResult: true,
	},
	{
		testDescription: "not matching map with string slice",
		requiredClaims: map[string]interface{}{
			"foo": map[string][]string{
				"foo": {"foobar"},
			},
		},
		tokenClaims: map[string]interface{}{
			"foo": map[string][]string{
				"foo": {"foo", "bar", "baz"},
			},
		},
		expectedResult: false,
	},
	{
		testDescription: "matching slice with map",
		requiredClaims: map[string]interface{}{
			"foo": []map[string]string{
				{"bar": "baz"},
			},
		},
		tokenClaims: map[string]interface{}{
			"foo": []map[string]string{
				{"bar": "baz"},
			},
		},
		expectedResult: true,
	},
	{
		testDescription: "not matching slice with map",
		requiredClaims: map[string]interface{}{
			"foo": []map[string]string{
				{"bar": "foobar"},
			},
		},
		tokenClaims: map[string]interface{}{
			"foo": []map[string]string{
				{"bar": "baz"},
			},
		},
		expectedResult: false,
	},
	{
		testDescription: "matching primitive types, slice and map",
		requiredClaims: map[string]interface{}{
			"foo": "bar",
			"bar": 1337,
			"baz": []string{"foo"},
			"oof": []map[string]string{
				{"bar": "baz"},
			},
		},
		tokenClaims: map[string]interface{}{
			"foo": "bar",
			"bar": 1337,
			"baz": []string{"foo"},
			"oof": []map[string]string{
				{"bar": "baz"},
			},
		},
		expectedResult: true,
	},
	{
		testDescription: "matching primitive types, slice and map where token contains multiple values",
		requiredClaims: map[string]interface{}{
			"foo": "bar",
			"bar": 1337,
			"baz": []string{"bar"},
			"oof": []map[string]string{
				{"bar": "baz"},
			},
		},
		tokenClaims: map[string]interface{}{
			"foo": "bar",
			"bar": 1337,
			"baz": []string{"foo", "bar", "baz"},
			"oof": []map[string]string{
				{"a": "b"},
				{"bar": "baz"},
				{"c": "d"},
			},
		},
		expectedResult: true,
	},
	{
		testDescription: "valid interface list in an interface map",
		requiredClaims: map[string]interface{}{
			"foo": map[string][]string{
				"bar": {"baz"},
			},
		},
		tokenClaims: map[string]interface{}{
			"foo": map[string]interface{}{
				"bar": []interface{}{
					"uno",
					"dos",
					"baz",
					"tres",
				},
			},
		},
		expectedResult: true,
	},
	{
		testDescription: "invalid interface list in an interface map",
		requiredClaims: map[string]interface{}{
			"foo": map[string][]string{
				"bar": {"baz"},
			},
		},
		tokenClaims: map[string]interface{}{
			"foo": map[string]interface{}{
				"bar": []interface{}{
					"uno",
					"dos",
					"tres",
				},
			},
		},
		expectedResult: false,
	},
}
