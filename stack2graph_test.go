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

// TestParsePackageName tests the ParsePackageName function with various inputs,
// including all previously defined test cases from TestExtractPackageNames.
func TestParsePackageName(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		// Test cases from the original TestParsePackageName
		{
			input:    "bitbucket.org/zetaactions/trinity/api/zivo/zivoconnect.NewZivoAPIHandler.func1",
			expected: "bitbucket.org/zetaactions/trinity/api/zivo/zivoconnect",
		},
		{
			input:    "bitbucket.org/zetaactions/trinity/api/zuri/zuriconnect.(*zuriAPIClient).GetRule",
			expected: "bitbucket.org/zetaactions/trinity/api/zuri/zuriconnect",
		},
		{
			input:    "bitbucket.org/zetaactions/trinity/kosmos/cmd.GrpcServer",
			expected: "bitbucket.org/zetaactions/trinity/kosmos/cmd",
		},
		{
			input:    "main.run.func2",
			expected: "main",
		},
		{
			input:    "net/http.(*ServeMux).ServeHTTP",
			expected: "net/http",
		},
		{
			input:    "created by net/http.",
			expected: "created by net/http",
		},
		{
			input:    "github.com/wricardo/stacktrace-to-graph.ReportStacktrace",
			expected: "github.com/wricardo/stacktrace-to-graph",
		},
		{
			input:    "golang.org/x/net/http2/h2c.h2cHandler.ServeHTTP",
			expected: "golang.org/x/net/http2/h2c",
		},
		{
			input:    "connectrpc.com/connect.NewUnaryHandler[...].func2",
			expected: "connectrpc.com/connect",
		},
		// {
		// 	input:    "some/unknownformatfunction",
		// 	expected: "some/unknownformatfunction",
		// },
		{
			input:    "NoSlashFunction",
			expected: "NoSlashFunction",
		},
		{
			input:    "just/a/path/with/no.function",
			expected: "just/a/path/with/no",
		},
		{
			input:    "bitbucket.org/zetaactions/trinity/api/zuri/zuriconnect.(*zuriAPIClient).GetRule",
			expected: "bitbucket.org/zetaactions/trinity/api/zuri/zuriconnect",
		},
		{
			input:    "bitbucket.org/zetaactions/trinity/kosmos/cmd.GrpcServer.DefaultClientInterceptors.NewInterceptor.func6.1",
			expected: "bitbucket.org/zetaactions/trinity/kosmos/cmd",
		},
		{
			input:    "bitbucket.org/zetaactions/trinity/kosmos/cmd.TemporalWorker",
			expected: "bitbucket.org/zetaactions/trinity/kosmos/cmd",
		},
		{
			input:    "bitbucket.org/zetaactions/trinity/kosmos/cmd.TemporalWorker.DefaultClientInterceptors.NewInterceptor.func2.1",
			expected: "bitbucket.org/zetaactions/trinity/kosmos/cmd",
		},
		{
			input:    "bitbucket.org/zetaactions/trinity/kosmos/router.Init",
			expected: "bitbucket.org/zetaactions/trinity/kosmos/router",
		},
		{
			input:    "bitbucket.org/zetaactions/trinity/kosmos/services.(*CampaignService).LoadZuriRules",
			expected: "bitbucket.org/zetaactions/trinity/kosmos/services",
		},
		{
			input:    "bitbucket.org/zetaactions/trinity/kosmos/services.(*CampaignService).LoadZuriRulesForCampaign",
			expected: "bitbucket.org/zetaactions/trinity/kosmos/services",
		},
		{
			input:    "bitbucket.org/zetaactions/trinity/kosmos/services.(*CampaignService).ReCache",
			expected: "bitbucket.org/zetaactions/trinity/kosmos/services",
		},
		{
			input:    "bitbucket.org/zetaactions/trinity/kosmos/services.NewCampaignService",
			expected: "bitbucket.org/zetaactions/trinity/kosmos/services",
		},
		{
			input:    "bitbucket.org/zetaactions/trinity/lib/ctxdata.FromContext",
			expected: "bitbucket.org/zetaactions/trinity/lib/ctxdata",
		},
		{
			input:    "bitbucket.org/zetaactions/trinity/lib/ctxdata.IsTest",
			expected: "bitbucket.org/zetaactions/trinity/lib/ctxdata",
		},
		{
			input:    "bitbucket.org/zetaactions/trinity/zivo/cmd.ApiServer.(*Middleware).Wrap.func3",
			expected: "bitbucket.org/zetaactions/trinity/zivo/cmd",
		},
		{
			input:    "bitbucket.org/zetaactions/trinity/zivo/cmd.ApiServer.NewInterceptor.func2.1",
			expected: "bitbucket.org/zetaactions/trinity/zivo/cmd",
		},
		{
			input:    "bitbucket.org/zetaactions/trinity/zivo/grpc.(*Handler).SendSms",
			expected: "bitbucket.org/zetaactions/trinity/zivo/grpc",
		},
		{
			input:    "bitbucket.org/zetaactions/trinity/zivo/services.(*SmsProviderFactory).BuildSmsProvider",
			expected: "bitbucket.org/zetaactions/trinity/zivo/services",
		},
		{
			input:    "bitbucket.org/zetaactions/trinity/zivo/services.(*bandwidthclient).SendSms",
			expected: "bitbucket.org/zetaactions/trinity/zivo/services",
		},
		{
			input:    "connectrpc.com/connect.(*Client[...]).CallUnary",
			expected: "connectrpc.com/connect",
		},
		{
			input:    "connectrpc.com/connect.(*Handler).ServeHTTP",
			expected: "connectrpc.com/connect",
		},
		{
			input:    "connectrpc.com/connect.NewClient[...].func2",
			expected: "connectrpc.com/connect",
		},
		{
			input:    "connectrpc.com/connect.NewUnaryHandler[...].func1",
			expected: "connectrpc.com/connect",
		},
		{
			input:    "connectrpc.com/connect.NewUnaryHandler[...].func2",
			expected: "connectrpc.com/connect",
		},
		{
			input:    "created by main.run in goroutine 1",
			expected: "created by main",
		},
		{
			input:    "github.com/wricardo/stacktrace-to-graph.(*StackToGraph).ReportStacktrace",
			expected: "github.com/wricardo/stacktrace-to-graph",
		},
		{
			input:    "github.com/wricardo/stacktrace-to-graph.ReportStacktrace",
			expected: "github.com/wricardo/stacktrace-to-graph",
		},
		{
			input:    "github.com/wricardo/stacktrace-to-graph.captureStackTrace",
			expected: "github.com/wricardo/stacktrace-to-graph",
		},
		{
			input:    "golang.org/x/net/http2/h2c.h2cHandler.ServeHTTP",
			expected: "golang.org/x/net/http2/h2c",
		},
		{
			input:    "main.run.func9",
			expected: "main",
		},
		{
			input:    "net/http.(*conn).serve",
			expected: "net/http",
		},
		{
			input:    "net/http.HandlerFunc.ServeHTTP",
			expected: "net/http",
		},
		{
			input:    "net/http.serverHandler.ServeHTTP",
			expected: "net/http",
		},
	}

	for _, tc := range testCases {
		result := ParsePackageName(tc.input)
		if result != tc.expected {
			t.Errorf("ParsePackageName(%q) = %q; expected %q", tc.input, result, tc.expected)
		}
	}
}
