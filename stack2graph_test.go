package stacktracetograph

import (
	"testing"
)

func TestCleanFunctionName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		// Standard Function Signatures
		{
			input:    "main.isError({0x104486073?, 0xc?})",
			expected: "main.isError",
		},
		{
			input:    "main.(*Person).SayHello(0x1400010aeb8)",
			expected: "main.(*Person).SayHello",
		},
		{
			input:    "main.(*Person).SayHello!@#(0x1400010aeb8)",
			expected: "main.(*Person).SayHello!@#",
		},
		// No Parentheses
		{
			input:    "main.isError",
			expected: "main.isError",
		},
		// Multiple Parentheses
		{
			input:    "main.isError(arg1, arg2(arg3))",
			expected: "main.isError",
		},
		// Nested Function Calls
		{
			input:    "main.funcA(funcB(arg))",
			expected: "main.funcA",
		},
		// Leading/Trailing Whitespaces
		{
			input:    "  main.isError(arg)  ",
			expected: "  main.isError",
		},
		// Empty String
		{
			input:    "",
			expected: "",
		},
		// Function Name with Special Characters
		{
			input:    "main.(*Person).SayHello!@#(0x1400010aeb8)",
			expected: "main.(*Person).SayHello!@#",
		},
		// Additional Test Cases
		{
			input:    "main.funcA  (arg)",
			expected: "main.funcA  ",
		},
		{
			input:    "(arg)",
			expected: "",
		},
		{
			input:    "main.func-A(arg)",
			expected: "main.func-A",
		},
		{
			input:    "main.123Func(arg)",
			expected: "main.123Func",
		},
		{
			input:    "main.你好函数(arg)",
			expected: "main.你好函数",
		},
		{
			input:    "main.funcA()(arg)",
			expected: "main.funcA()",
		},
	}

	for _, test := range tests {
		result := cleanFunctionName(test.input)
		if result != test.expected {
			t.Errorf("cleanFunctionName(%q) = %q; want %q", test.input, result, test.expected)
		}
	}
}
