package templates

import "fmt"

func GetEditorConfig(template string) string {
	switch template {
	case "go":
		return `root = true

[*]
charset = utf-8
end_of_line = lf
insert_final_newline = true
trim_trailing_whitespace = true

[*.go]
indent_style = tab
indent_size = 4

[*.{mod,sum}]
indent_style = tab
indent_size = 4

[*.md]
trim_trailing_whitespace = false

[Makefile]
indent_style = tab
`
	case "python":
		return `root = true

[*]
charset = utf-8
end_of_line = lf
insert_final_newline = true
trim_trailing_whitespace = true

[*.py]
indent_style = space
indent_size = 4

[*.{cfg,ini,toml}]
indent_style = space
indent_size = 4

[*.md]
trim_trailing_whitespace = false
`
	case "node", "typescript":
		return `root = true

[*]
charset = utf-8
end_of_line = lf
insert_final_newline = true
trim_trailing_whitespace = true

[*.{js,ts,jsx,tsx}]
indent_style = space
indent_size = 2

[*.{json,yml,yaml}]
indent_style = space
indent_size = 2

[*.md]
trim_trailing_whitespace = false
`
	case "rust":
		return `root = true

[*]
charset = utf-8
end_of_line = lf
insert_final_newline = true
trim_trailing_whitespace = true

[*.rs]
indent_style = space
indent_size = 4

[Cargo.toml]
indent_style = space
indent_size = 4

[*.md]
trim_trailing_whitespace = false
`
	case "java", "kotlin":
		return `root = true

[*]
charset = utf-8
end_of_line = lf
insert_final_newline = true
trim_trailing_whitespace = true

[*.{java,kt}]
indent_style = space
indent_size = 4

[*.{gradle,groovy}]
indent_style = space
indent_size = 4

[*.xml]
indent_style = space
indent_size = 4

[*.md]
trim_trailing_whitespace = false
`
	case "csharp":
		return `root = true

[*]
charset = utf-8
end_of_line = lf
insert_final_newline = true
trim_trailing_whitespace = true

[*.{cs,vb,fs}]
indent_style = space
indent_size = 4

[*.{csproj,sln}]
indent_style = space
indent_size = 2

[*.md]
trim_trailing_whitespace = false
`
	case "php":
		return `root = true

[*]
charset = utf-8
end_of_line = lf
insert_final_newline = true
trim_trailing_whitespace = true

[*.php]
indent_style = space
indent_size = 4

[*.md]
trim_trailing_whitespace = false
`
	case "ruby":
		return `root = true

[*]
charset = utf-8
end_of_line = lf
insert_final_newline = true
trim_trailing_whitespace = true

[*.{rb,rake}]
indent_style = space
indent_size = 2

[*.md]
trim_trailing_whitespace = false
`
	case "swift":
		return `root = true

[*]
charset = utf-8
end_of_line = lf
insert_final_newline = true
trim_trailing_whitespace = true

[*.swift]
indent_style = space
indent_size = 4

[*.md]
trim_trailing_whitespace = false
`
	case "dart":
		return `root = true

[*]
charset = utf-8
end_of_line = lf
insert_final_newline = true
trim_trailing_whitespace = true

[*.dart]
indent_style = space
indent_size = 2

[*.yaml]
indent_style = space
indent_size = 2

[*.md]
trim_trailing_whitespace = false
`
	case "c-cpp":
		return `root = true

[*]
charset = utf-8
end_of_line = lf
insert_final_newline = true
trim_trailing_whitespace = true

[*.{c,h,cpp,hpp}]
indent_style = space
indent_size = 4

[CMakeLists.txt]
indent_style = space
indent_size = 4

[*.md]
trim_trailing_whitespace = false
`
	default:
		return `root = true

[*]
charset = utf-8
end_of_line = lf
insert_final_newline = true
trim_trailing_whitespace = true

[*.{go,py,js,ts,rb,java,rs,cs,php,swift,kt,dart,c,cpp}]
indent_style = space
indent_size = 4

[*.{json,yml,yaml,toml}]
indent_style = space
indent_size = 2

[Makefile]
indent_style = tab

[*.md]
trim_trailing_whitespace = false
`
	}
}

func GetCodeowners(content string) string {
	if content == "" {
		content = "* @owner"
	}
	return fmt.Sprintf(`# Code Owners for this repository
# See: https://docs.github.com/en/repositories/managing-your-repositorys-settings-and-features/customizing-your-repository/about-code-owners

%s
`, content)
}

func GetSecurityMD(projectName, contact string) string {
	if contact == "" {
		contact = "security@example.com"
	}
	return fmt.Sprintf(`# Security Policy

## Supported Versions

| Version | Supported          |
| ------- | ------------------ |
| latest  | :white_check_mark: |

## Reporting a Vulnerability

If you discover a security vulnerability within %s, please send an e-mail to %s. All security vulnerabilities will be promptly addressed.

**Please do not report security vulnerabilities through public GitHub issues.**

## Response Process

1. Acknowledgment of receipt within 48 hours
2. Assessment of the vulnerability within 1 week
3. Resolution and disclosure timeline communicated to reporter

## Security Best Practices

- All dependencies are regularly updated
- CI/CD pipelines include security scanning
- Sensitive data is never committed to the repository

Thank you for helping keep %s and its users safe.
`, projectName, contact, projectName)
}
