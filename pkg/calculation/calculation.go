package calculation

import (
	"errors"
	"fmt"
	"strconv"
)

func Calc(expression string) (float64, error) {
	result, err := evaluateExpression(expression)
	if err != nil {
		return 0, err
	}
	return result, nil
}

func evaluateExpression(expression string) (float64, error) { // Проверяем, что выражение не пустое
	if len(expression) == 0 {
		return 0, errors.New("empty expression")
	}

	//преобразуем строку в обратную польскую запись
	postfixExpression, err := infixToPostfix(expression)
	if err != nil {
		return 0, err
	}

	//вычисляем значение выражения в обратной польской записи
	stack := make([]float64, 0)
	for _, token := range postfixExpression {
		if isOperator(token) {
			if len(stack) < 2 {
				return 0, errors.New("invalid expression")
			}
			operand2 := stack[len(stack)-1]
			operand1 := stack[len(stack)-2]
			stack = stack[:len(stack)-2]
			result := performOperation(operand1, operand2, token)
			stack = append(stack, result)
		} else {
			value, err := strconv.ParseFloat(token, 64)
			if err != nil {
				return 0, errors.New("invalid expression")
			}
			stack = append(stack, value)
		}
	}

	if len(stack) != 1 {
		return 0, errors.New("invalid expression")
	}

	return stack[0], nil
}

func infixToPostfix(expression string) ([]string, error) {
	var result []string
	var operators []string
	tokens := tokenize(expression)
	for _, token := range tokens {
		if isNumber(token) {
			result = append(result, token)
		} else if token == "(" {
			operators = append(operators, token)
		} else if token == ")" {
			for len(operators) > 0 && operators[len(operators)-1] != "(" {
				result = append(result, operators[len(operators)-1])
				operators = operators[:len(operators)-1]
			}
			if len(operators) == 0 {
				return nil, errors.New("unmatched parentheses")
			}
			operators = operators[:len(operators)-1]
		} else if isOperator(token) {
			for len(operators) > 0 && precedence(operators[len(operators)-1]) >= precedence(token) {
				result = append(result, operators[len(operators)-1])
				operators = operators[:len(operators)-1]
			}
			operators = append(operators, token)
		} else {
			return nil, errors.New("invalid token")
		}
	}

	for len(operators) > 0 {
		if operators[len(operators)-1] == "(" {
			return nil, errors.New("unmatched parentheses")
		}
		result = append(result, operators[len(operators)-1])
		operators = operators[:len(operators)-1]
	}

	return result, nil
}

func tokenize(expression string) []string {
	var tokens []string
	token := ""
	for _, char := range expression {
		if isOperator(string(char)) || string(char) == "(" || string(char) == ")" {
			if token != "" {
				tokens = append(tokens, token)
				token = ""
			}
			tokens = append(tokens, string(char))
		} else if string(char) == " " {
			continue
		} else {
			token += string(char)
		}
	}

	if token != "" {
		tokens = append(tokens, token)
	}

	return tokens
}

func isNumber(token string) bool {
	_, err := strconv.ParseFloat(token, 64)
	return err == nil
}

func isOperator(token string) bool {
	return token == "+" || token == "-" || token == "*" || token == "/"
}

func precedence(operator string) int {
	switch operator {
	case "+", "-":
		return 1
	case "*", "/":
		return 2
	default:
		return 0
	}
}

func performOperation(operand1, operand2 float64, operator string) float64 {
	switch operator {
	case "+":
		return operand1 + operand2
	case "-":
		return operand1 - operand2
	case "*":
		return operand1 * operand2
	case "/":
		return operand1 / operand2
	default:
		return 0
	}
}
func main() {
	var vhod string
	fmt.Scanln(&vhod)
	result, err := Calc(vhod)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%.2f\n", result)
}
