package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Set up logging to stderr so it doesn't interfere with stdio communication
	log.SetOutput(os.Stderr)
	log.SetPrefix("[Calculator MCP] ")
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	log.Println("Starting Calculator MCP Server v1.0.0")

	// Create a new MCP server
	s := server.NewMCPServer(
		"Calculator MCP Server",
		"1.0.0",
		server.WithToolCapabilities(false),
		server.WithRecovery(),
	)

	// Add basic arithmetic tools
	log.Println("Registering calculator tools...")
	addAdditionTool(s)
	addSubtractionTool(s)
	addMultiplicationTool(s)
	addDivisionTool(s)
	addPowerTool(s)
	addSquareRootTool(s)
	log.Println("All tools registered successfully")

	// Start the stdio server
	log.Println("Starting stdio server...")
	if err := server.ServeStdio(s); err != nil {
		log.Printf("Server error: %v\n", err)
		fmt.Printf("Server error: %v\n", err)
	}
}

// addAdditionTool adds the addition tool to the server
func addAdditionTool(s *server.MCPServer) {
	addTool := mcp.NewTool("add",
		mcp.WithDescription("Add two numbers together"),
		mcp.WithNumber("a",
			mcp.Required(),
			mcp.Description("First number to add"),
		),
		mcp.WithNumber("b",
			mcp.Required(),
			mcp.Description("Second number to add"),
		),
	)

	s.AddTool(addTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		log.Printf("ADD tool called")
		
		a, err := request.RequireFloat("a")
		if err != nil {
			log.Printf("ADD tool error - invalid parameter 'a': %v", err)
			return mcp.NewToolResultError(err.Error()), nil
		}

		b, err := request.RequireFloat("b")
		if err != nil {
			log.Printf("ADD tool error - invalid parameter 'b': %v", err)
			return mcp.NewToolResultError(err.Error()), nil
		}

		result := a + b
		log.Printf("ADD tool: %.6g + %.6g = %.6g", a, b, result)
		return mcp.NewToolResultText(fmt.Sprintf("%.6g", result)), nil
	})
}

// addSubtractionTool adds the subtraction tool to the server
func addSubtractionTool(s *server.MCPServer) {
	subtractTool := mcp.NewTool("subtract",
		mcp.WithDescription("Subtract the second number from the first number"),
		mcp.WithNumber("a",
			mcp.Required(),
			mcp.Description("Number to subtract from"),
		),
		mcp.WithNumber("b",
			mcp.Required(),
			mcp.Description("Number to subtract"),
		),
	)

	s.AddTool(subtractTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		log.Printf("SUBTRACT tool called")
		
		a, err := request.RequireFloat("a")
		if err != nil {
			log.Printf("SUBTRACT tool error - invalid parameter 'a': %v", err)
			return mcp.NewToolResultError(err.Error()), nil
		}

		b, err := request.RequireFloat("b")
		if err != nil {
			log.Printf("SUBTRACT tool error - invalid parameter 'b': %v", err)
			return mcp.NewToolResultError(err.Error()), nil
		}

		result := a - b
		log.Printf("SUBTRACT tool: %.6g - %.6g = %.6g", a, b, result)
		return mcp.NewToolResultText(fmt.Sprintf("%.6g", result)), nil
	})
}

// addMultiplicationTool adds the multiplication tool to the server
func addMultiplicationTool(s *server.MCPServer) {
	multiplyTool := mcp.NewTool("multiply",
		mcp.WithDescription("Multiply two numbers together"),
		mcp.WithNumber("a",
			mcp.Required(),
			mcp.Description("First number to multiply"),
		),
		mcp.WithNumber("b",
			mcp.Required(),
			mcp.Description("Second number to multiply"),
		),
	)

	s.AddTool(multiplyTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		log.Printf("MULTIPLY tool called")
		
		a, err := request.RequireFloat("a")
		if err != nil {
			log.Printf("MULTIPLY tool error - invalid parameter 'a': %v", err)
			return mcp.NewToolResultError(err.Error()), nil
		}

		b, err := request.RequireFloat("b")
		if err != nil {
			log.Printf("MULTIPLY tool error - invalid parameter 'b': %v", err)
			return mcp.NewToolResultError(err.Error()), nil
		}

		result := a * b
		log.Printf("MULTIPLY tool: %.6g * %.6g = %.6g", a, b, result)
		return mcp.NewToolResultText(fmt.Sprintf("%.6g", result)), nil
	})
}

// addDivisionTool adds the division tool to the server
func addDivisionTool(s *server.MCPServer) {
	divideTool := mcp.NewTool("divide",
		mcp.WithDescription("Divide the first number by the second number"),
		mcp.WithNumber("a",
			mcp.Required(),
			mcp.Description("Number to be divided (dividend)"),
		),
		mcp.WithNumber("b",
			mcp.Required(),
			mcp.Description("Number to divide by (divisor)"),
		),
	)

	s.AddTool(divideTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		log.Printf("DIVIDE tool called")
		
		a, err := request.RequireFloat("a")
		if err != nil {
			log.Printf("DIVIDE tool error - invalid parameter 'a': %v", err)
			return mcp.NewToolResultError(err.Error()), nil
		}

		b, err := request.RequireFloat("b")
		if err != nil {
			log.Printf("DIVIDE tool error - invalid parameter 'b': %v", err)
			return mcp.NewToolResultError(err.Error()), nil
		}

		if b == 0 {
			log.Printf("DIVIDE tool error - division by zero: %.6g / %.6g", a, b)
			return mcp.NewToolResultError("cannot divide by zero"), nil
		}

		result := a / b
		log.Printf("DIVIDE tool: %.6g / %.6g = %.6g", a, b, result)
		return mcp.NewToolResultText(fmt.Sprintf("%.6g", result)), nil
	})
}

// addPowerTool adds the power/exponentiation tool to the server
func addPowerTool(s *server.MCPServer) {
	powerTool := mcp.NewTool("power",
		mcp.WithDescription("Raise the first number to the power of the second number"),
		mcp.WithNumber("base",
			mcp.Required(),
			mcp.Description("Base number"),
		),
		mcp.WithNumber("exponent",
			mcp.Required(),
			mcp.Description("Exponent (power to raise the base to)"),
		),
	)

	s.AddTool(powerTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		log.Printf("POWER tool called")
		
		base, err := request.RequireFloat("base")
		if err != nil {
			log.Printf("POWER tool error - invalid parameter 'base': %v", err)
			return mcp.NewToolResultError(err.Error()), nil
		}

		exponent, err := request.RequireFloat("exponent")
		if err != nil {
			log.Printf("POWER tool error - invalid parameter 'exponent': %v", err)
			return mcp.NewToolResultError(err.Error()), nil
		}

		result := math.Pow(base, exponent)
		
		// Check for invalid results
		if math.IsNaN(result) {
			log.Printf("POWER tool error - result is NaN: %.6g ^ %.6g", base, exponent)
			return mcp.NewToolResultError("result is not a number (NaN)"), nil
		}
		if math.IsInf(result, 0) {
			log.Printf("POWER tool error - result is infinite: %.6g ^ %.6g", base, exponent)
			return mcp.NewToolResultError("result is infinite"), nil
		}

		log.Printf("POWER tool: %.6g ^ %.6g = %.6g", base, exponent, result)
		return mcp.NewToolResultText(fmt.Sprintf("%.6g", result)), nil
	})
}

// addSquareRootTool adds the square root tool to the server
func addSquareRootTool(s *server.MCPServer) {
	sqrtTool := mcp.NewTool("sqrt",
		mcp.WithDescription("Calculate the square root of a number"),
		mcp.WithNumber("number",
			mcp.Required(),
			mcp.Description("Number to calculate the square root of"),
		),
	)

	s.AddTool(sqrtTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		log.Printf("SQRT tool called")
		
		number, err := request.RequireFloat("number")
		if err != nil {
			log.Printf("SQRT tool error - invalid parameter 'number': %v", err)
			return mcp.NewToolResultError(err.Error()), nil
		}

		if number < 0 {
			log.Printf("SQRT tool error - negative number: %.6g", number)
			return mcp.NewToolResultError("cannot calculate square root of negative number"), nil
		}

		result := math.Sqrt(number)
		log.Printf("SQRT tool: √%.6g = %.6g", number, result)
		return mcp.NewToolResultText(fmt.Sprintf("%.6g", result)), nil
	})
}