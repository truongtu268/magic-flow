package codegen

import (
	"fmt"
	"path/filepath"
	"strings"

	"magic-flow/v2/pkg/models"
)

// JavaHandler implements LanguageHandler for Java code generation
type JavaHandler struct {
	templateManager *TemplateManager
}

// NewJavaHandler creates a new Java language handler
func NewJavaHandler(templateManager *TemplateManager) *JavaHandler {
	return &JavaHandler{
		templateManager: templateManager,
	}
}

// Generate generates Java code for a workflow
func (h *JavaHandler) Generate(workflow *models.Workflow, request *GenerationRequest, templateData *TemplateData) ([]GeneratedFile, error) {
	var files []GeneratedFile

	// Generate client file
	clientFile, err := h.generateClientFile(templateData)
	if err != nil {
		return nil, fmt.Errorf("failed to generate client file: %w", err)
	}
	files = append(files, clientFile)

	// Generate models files
	modelFiles, err := h.generateModelFiles(templateData)
	if err != nil {
		return nil, fmt.Errorf("failed to generate model files: %w", err)
	}
	files = append(files, modelFiles...)

	// Generate exception files
	exceptionFiles, err := h.generateExceptionFiles(templateData)
	if err != nil {
		return nil, fmt.Errorf("failed to generate exception files: %w", err)
	}
	files = append(files, exceptionFiles...)

	// Generate configuration file
	configFile, err := h.generateConfigFile(templateData)
	if err != nil {
		return nil, fmt.Errorf("failed to generate config file: %w", err)
	}
	files = append(files, configFile)

	// Generate test file if requested
	if request.IncludeTests {
		testFile, err := h.generateTestFile(templateData)
		if err != nil {
			return nil, fmt.Errorf("failed to generate test file: %w", err)
		}
		files = append(files, testFile)
	}

	// Generate pom.xml file
	pomFile, err := h.generatePomFile(templateData)
	if err != nil {
		return nil, fmt.Errorf("failed to generate pom.xml file: %w", err)
	}
	files = append(files, pomFile)

	// Generate gradle build file
	gradleFile, err := h.generateGradleFile(templateData)
	if err != nil {
		return nil, fmt.Errorf("failed to generate build.gradle file: %w", err)
	}
	files = append(files, gradleFile)

	// Generate README file
	readmeFile, err := h.generateReadmeFile(templateData)
	if err != nil {
		return nil, fmt.Errorf("failed to generate README file: %w", err)
	}
	files = append(files, readmeFile)

	return files, nil
}

// ValidateRequest validates Java-specific generation request
func (h *JavaHandler) ValidateRequest(request *GenerationRequest) error {
	if request.PackageName == "" {
		request.PackageName = h.GetDefaultPackageName()
	}

	// Validate package name format (Java package naming rules)
	if !isValidJavaPackageName(request.PackageName) {
		return fmt.Errorf("invalid Java package name: %s", request.PackageName)
	}

	return nil
}

// PrepareTemplateData prepares template data for Java code generation
func (h *JavaHandler) PrepareTemplateData(workflow *models.Workflow, request *GenerationRequest) (*TemplateData, error) {
	packageName := request.PackageName
	if packageName == "" {
		packageName = h.GetDefaultPackageName()
	}

	className := ToPascalCase(workflow.Name) + "Client"

	// Extract methods from workflow steps
	methods := ExtractStepMethods(workflow)

	// Generate imports
	imports := h.generateImports(workflow, request)

	// Generate models
	models := h.generateModels(workflow)

	templateData := &TemplateData{
		Workflow:    workflow,
		PackageName: packageName,
		ClassName:   className,
		Imports:     imports,
		Methods:     methods,
		Models:      models,
		Options:     request.Options,
		GeneratedAt: workflow.CreatedAt,
	}

	return templateData, nil
}

// GetFileExtension returns the file extension for Java files
func (h *JavaHandler) GetFileExtension() string {
	return ".java"
}

// GetDefaultPackageName returns the default package name for Java
func (h *JavaHandler) GetDefaultPackageName() string {
	return "com.magicflow.client"
}

// generateClientFile generates the main client file
func (h *JavaHandler) generateClientFile(data *TemplateData) (GeneratedFile, error) {
	template, err := h.templateManager.GetTemplate("java", "client")
	if err != nil {
		return GeneratedFile{}, err
	}

	content, err := RenderTemplate(template, data)
	if err != nil {
		return GeneratedFile{}, err
	}

	packagePath := strings.ReplaceAll(data.PackageName, ".", "/")
	filePath := filepath.Join("src", "main", "java", packagePath, data.ClassName+".java")

	return GeneratedFile{
		Path:     filePath,
		Content:  content,
		Language: "java",
		Type:     "client",
	}, nil
}

// generateModelFiles generates model files
func (h *JavaHandler) generateModelFiles(data *TemplateData) ([]GeneratedFile, error) {
	var files []GeneratedFile
	packagePath := strings.ReplaceAll(data.PackageName, ".", "/")

	// Generate base model classes
	baseModels := []string{"ExecutionResult", "ExecutionStatus", "StepStatus", "WorkflowInput", "WorkflowOutput"}

	for _, modelName := range baseModels {
		content := h.generateBaseModelContent(data, modelName)
		filePath := filepath.Join("src", "main", "java", packagePath, "models", modelName+".java")

		files = append(files, GeneratedFile{
			Path:     filePath,
			Content:  content,
			Language: "java",
			Type:     "model",
		})
	}

	// Generate workflow-specific models
	for _, model := range data.Models {
		content := h.generateModelContent(data, model)
		filePath := filepath.Join("src", "main", "java", packagePath, "models", model.Name+".java")

		files = append(files, GeneratedFile{
			Path:     filePath,
			Content:  content,
			Language: "java",
			Type:     "model",
		})
	}

	return files, nil
}

// generateExceptionFiles generates exception files
func (h *JavaHandler) generateExceptionFiles(data *TemplateData) ([]GeneratedFile, error) {
	var files []GeneratedFile
	packagePath := strings.ReplaceAll(data.PackageName, ".", "/")

	exceptions := []string{
		"MagicFlowException",
		"ApiException",
		"AuthenticationException",
		"ValidationException",
		"ExecutionException",
		"TimeoutException",
		"NetworkException",
	}

	for _, exceptionName := range exceptions {
		content := h.generateExceptionContent(data, exceptionName)
		filePath := filepath.Join("src", "main", "java", packagePath, "exceptions", exceptionName+".java")

		files = append(files, GeneratedFile{
			Path:     filePath,
			Content:  content,
			Language: "java",
			Type:     "exception",
		})
	}

	return files, nil
}

// generateConfigFile generates the configuration file
func (h *JavaHandler) generateConfigFile(data *TemplateData) (GeneratedFile, error) {
	packagePath := strings.ReplaceAll(data.PackageName, ".", "/")

	content := fmt.Sprintf(`package %s.config;

import java.time.Duration;

/**
 * Configuration class for %s Client
 * Generated at: %s
 */
public class ClientConfig {
    public static final String WORKFLOW_ID = "%s";
    public static final String WORKFLOW_NAME = "%s";
    public static final String DEFAULT_BASE_URL = "http://localhost:8080";
    public static final Duration DEFAULT_TIMEOUT = Duration.ofSeconds(30);
    public static final int DEFAULT_RETRY_ATTEMPTS = 3;
    public static final Duration DEFAULT_RETRY_DELAY = Duration.ofSeconds(1);
    
    private String baseUrl;
    private String apiKey;
    private Duration timeout;
    private int retryAttempts;
    private Duration retryDelay;
    private boolean debug;
    
    public ClientConfig() {
        this.baseUrl = DEFAULT_BASE_URL;
        this.timeout = DEFAULT_TIMEOUT;
        this.retryAttempts = DEFAULT_RETRY_ATTEMPTS;
        this.retryDelay = DEFAULT_RETRY_DELAY;
        this.debug = false;
    }
    
    public ClientConfig(String baseUrl, String apiKey) {
        this();
        this.baseUrl = baseUrl;
        this.apiKey = apiKey;
    }
    
    // Getters and setters
    public String getBaseUrl() { return baseUrl; }
    public void setBaseUrl(String baseUrl) { this.baseUrl = baseUrl; }
    
    public String getApiKey() { return apiKey; }
    public void setApiKey(String apiKey) { this.apiKey = apiKey; }
    
    public Duration getTimeout() { return timeout; }
    public void setTimeout(Duration timeout) { this.timeout = timeout; }
    
    public int getRetryAttempts() { return retryAttempts; }
    public void setRetryAttempts(int retryAttempts) { this.retryAttempts = retryAttempts; }
    
    public Duration getRetryDelay() { return retryDelay; }
    public void setRetryDelay(Duration retryDelay) { this.retryDelay = retryDelay; }
    
    public boolean isDebug() { return debug; }
    public void setDebug(boolean debug) { this.debug = debug; }
}
`,
		data.PackageName,
		data.Workflow.Name,
		data.GeneratedAt.Format("2006-01-02 15:04:05"),
		data.Workflow.ID.String(),
		data.Workflow.Name,
	)

	packagePath = strings.ReplaceAll(data.PackageName, ".", "/")
	filePath := filepath.Join("src", "main", "java", packagePath, "config", "ClientConfig.java")

	return GeneratedFile{
		Path:     filePath,
		Content:  content,
		Language: "java",
		Type:     "config",
	}, nil
}

// generateTestFile generates the test file
func (h *JavaHandler) generateTestFile(data *TemplateData) (GeneratedFile, error) {
	template, err := h.templateManager.GetTemplate("java", "test")
	if err != nil {
		return GeneratedFile{}, err
	}

	content, err := RenderTemplate(template, data)
	if err != nil {
		return GeneratedFile{}, err
	}

	packagePath := strings.ReplaceAll(data.PackageName, ".", "/")
	filePath := filepath.Join("src", "test", "java", packagePath, data.ClassName+"Test.java")

	return GeneratedFile{
		Path:     filePath,
		Content:  content,
		Language: "java",
		Type:     "test",
	}, nil
}

// generatePomFile generates the pom.xml file
func (h *JavaHandler) generatePomFile(data *TemplateData) (GeneratedFile, error) {
	version := "1.0.0"
	if data.Options != nil {
		if v, ok := data.Options["version"].(string); ok && v != "" {
			version = v
		}
	}

	groupId := "com.magicflow"
	if data.Options != nil {
		if g, ok := data.Options["group_id"].(string); ok && g != "" {
			groupId = g
		}
	}

	artifactId := ToSnakeCase(data.Workflow.Name) + "-client"
	if data.Options != nil {
		if a, ok := data.Options["artifact_id"].(string); ok && a != "" {
			artifactId = a
		}
	}

	content := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<project xmlns="http://maven.apache.org/POM/4.0.0"
         xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
         xsi:schemaLocation="http://maven.apache.org/POM/4.0.0 
         http://maven.apache.org/xsd/maven-4.0.0.xsd">
    <modelVersion>4.0.0</modelVersion>

    <groupId>%s</groupId>
    <artifactId>%s</artifactId>
    <version>%s</version>
    <packaging>jar</packaging>

    <name>%s Java Client</name>
    <description>Java client library for %s workflow</description>
    <url>https://github.com/your-org/your-repo</url>

    <licenses>
        <license>
            <name>MIT License</name>
            <url>https://opensource.org/licenses/MIT</url>
        </license>
    </licenses>

    <developers>
        <developer>
            <name>Magic Flow</name>
            <email>contact@magicflow.dev</email>
            <organization>Magic Flow</organization>
            <organizationUrl>https://magicflow.dev</organizationUrl>
        </developer>
    </developers>

    <scm>
        <connection>scm:git:git://github.com/your-org/your-repo.git</connection>
        <developerConnection>scm:git:ssh://github.com:your-org/your-repo.git</developerConnection>
        <url>https://github.com/your-org/your-repo/tree/main</url>
    </scm>

    <properties>
        <maven.compiler.source>11</maven.compiler.source>
        <maven.compiler.target>11</maven.compiler.target>
        <project.build.sourceEncoding>UTF-8</project.build.sourceEncoding>
        <jackson.version>2.15.2</jackson.version>
        <okhttp.version>4.11.0</okhttp.version>
        <junit.version>5.10.0</junit.version>
        <mockito.version>5.5.0</mockito.version>
    </properties>

    <dependencies>
        <!-- HTTP Client -->
        <dependency>
            <groupId>com.squareup.okhttp3</groupId>
            <artifactId>okhttp</artifactId>
            <version>${okhttp.version}</version>
        </dependency>
        
        <!-- JSON Processing -->
        <dependency>
            <groupId>com.fasterxml.jackson.core</groupId>
            <artifactId>jackson-databind</artifactId>
            <version>${jackson.version}</version>
        </dependency>
        <dependency>
            <groupId>com.fasterxml.jackson.datatype</groupId>
            <artifactId>jackson-datatype-jsr310</artifactId>
            <version>${jackson.version}</version>
        </dependency>
        
        <!-- Logging -->
        <dependency>
            <groupId>org.slf4j</groupId>
            <artifactId>slf4j-api</artifactId>
            <version>2.0.7</version>
        </dependency>
        
        <!-- Test Dependencies -->
        <dependency>
            <groupId>org.junit.jupiter</groupId>
            <artifactId>junit-jupiter</artifactId>
            <version>${junit.version}</version>
            <scope>test</scope>
        </dependency>
        <dependency>
            <groupId>org.mockito</groupId>
            <artifactId>mockito-core</artifactId>
            <version>${mockito.version}</version>
            <scope>test</scope>
        </dependency>
        <dependency>
            <groupId>org.mockito</groupId>
            <artifactId>mockito-junit-jupiter</artifactId>
            <version>${mockito.version}</version>
            <scope>test</scope>
        </dependency>
        <dependency>
            <groupId>com.squareup.okhttp3</groupId>
            <artifactId>mockwebserver</artifactId>
            <version>${okhttp.version}</version>
            <scope>test</scope>
        </dependency>
    </dependencies>

    <build>
        <plugins>
            <plugin>
                <groupId>org.apache.maven.plugins</groupId>
                <artifactId>maven-compiler-plugin</artifactId>
                <version>3.11.0</version>
                <configuration>
                    <source>11</source>
                    <target>11</target>
                </configuration>
            </plugin>
            
            <plugin>
                <groupId>org.apache.maven.plugins</groupId>
                <artifactId>maven-surefire-plugin</artifactId>
                <version>3.1.2</version>
            </plugin>
            
            <plugin>
                <groupId>org.apache.maven.plugins</groupId>
                <artifactId>maven-source-plugin</artifactId>
                <version>3.3.0</version>
                <executions>
                    <execution>
                        <id>attach-sources</id>
                        <goals>
                            <goal>jar</goal>
                        </goals>
                    </execution>
                </executions>
            </plugin>
            
            <plugin>
                <groupId>org.apache.maven.plugins</groupId>
                <artifactId>maven-javadoc-plugin</artifactId>
                <version>3.5.0</version>
                <executions>
                    <execution>
                        <id>attach-javadocs</id>
                        <goals>
                            <goal>jar</goal>
                        </goals>
                    </execution>
                </executions>
            </plugin>
            
            <plugin>
                <groupId>org.jacoco</groupId>
                <artifactId>jacoco-maven-plugin</artifactId>
                <version>0.8.10</version>
                <executions>
                    <execution>
                        <goals>
                            <goal>prepare-agent</goal>
                        </goals>
                    </execution>
                    <execution>
                        <id>report</id>
                        <phase>test</phase>
                        <goals>
                            <goal>report</goal>
                        </goals>
                    </execution>
                </executions>
            </plugin>
        </plugins>
    </build>
</project>
`,
		groupId,
		artifactId,
		version,
		data.Workflow.Name,
		data.Workflow.Name,
	)

	return GeneratedFile{
		Path:     "pom.xml",
		Content:  content,
		Language: "xml",
		Type:     "config",
	}, nil
}

// generateGradleFile generates the build.gradle file
func (h *JavaHandler) generateGradleFile(data *TemplateData) (GeneratedFile, error) {
	version := "1.0.0"
	if data.Options != nil {
		if v, ok := data.Options["version"].(string); ok && v != "" {
			version = v
		}
	}

	content := fmt.Sprintf(`plugins {
    id 'java-library'
    id 'maven-publish'
    id 'signing'
    id 'jacoco'
}

group = 'com.magicflow'
version = '%s'
java.sourceCompatibility = JavaVersion.VERSION_11

repositories {
    mavenCentral()
}

dependencies {
    // HTTP Client
    implementation 'com.squareup.okhttp3:okhttp:4.11.0'
    
    // JSON Processing
    implementation 'com.fasterxml.jackson.core:jackson-databind:2.15.2'
    implementation 'com.fasterxml.jackson.datatype:jackson-datatype-jsr310:2.15.2'
    
    // Logging
    implementation 'org.slf4j:slf4j-api:2.0.7'
    
    // Test Dependencies
    testImplementation 'org.junit.jupiter:junit-jupiter:5.10.0'
    testImplementation 'org.mockito:mockito-core:5.5.0'
    testImplementation 'org.mockito:mockito-junit-jupiter:5.5.0'
    testImplementation 'com.squareup.okhttp3:mockwebserver:4.11.0'
}

tasks.named('test') {
    useJUnitPlatform()
    finalizedBy jacocoTestReport
}

jacocoTestReport {
    dependsOn test
    reports {
        xml.required = true
        html.required = true
    }
}

java {
    withJavadocJar()
    withSourcesJar()
}

publishing {
    publications {
        maven(MavenPublication) {
            from components.java
            
            pom {
                name = '%s Java Client'
                description = 'Java client library for %s workflow'
                url = 'https://github.com/your-org/your-repo'
                
                licenses {
                    license {
                        name = 'MIT License'
                        url = 'https://opensource.org/licenses/MIT'
                    }
                }
                
                developers {
                    developer {
                        name = 'Magic Flow'
                        email = 'contact@magicflow.dev'
                        organization = 'Magic Flow'
                        organizationUrl = 'https://magicflow.dev'
                    }
                }
                
                scm {
                    connection = 'scm:git:git://github.com/your-org/your-repo.git'
                    developerConnection = 'scm:git:ssh://github.com:your-org/your-repo.git'
                    url = 'https://github.com/your-org/your-repo/tree/main'
                }
            }
        }
    }
}

signing {
    sign publishing.publications.maven
}

tasks.withType(JavaCompile) {
    options.encoding = 'UTF-8'
}

tasks.withType(Javadoc) {
    options.encoding = 'UTF-8'
}
`,
		version,
		data.Workflow.Name,
		data.Workflow.Name,
	)

	return GeneratedFile{
		Path:     "build.gradle",
		Content:  content,
		Language: "gradle",
		Type:     "config",
	}, nil
}

// generateReadmeFile generates the README file
func (h *JavaHandler) generateReadmeFile(data *TemplateData) (GeneratedFile, error) {
	groupId := "com.magicflow"
	if data.Options != nil {
		if g, ok := data.Options["group_id"].(string); ok && g != "" {
			groupId = g
		}
	}

	artifactId := ToSnakeCase(data.Workflow.Name) + "-client"
	if data.Options != nil {
		if a, ok := data.Options["artifact_id"].(string); ok && a != "" {
			artifactId = a
		}
	}

	version := "1.0.0"
	if data.Options != nil {
		if v, ok := data.Options["version"].(string); ok && v != "" {
			version = v
		}
	}

	content := fmt.Sprintf(`# %s Java Client

Generated Java client library for the %s workflow.

## Installation

### Maven

Add the following dependency to your ` + "`pom.xml`" + `:

` + "```xml" + `
<dependency>
    <groupId>%s</groupId>
    <artifactId>%s</artifactId>
    <version>%s</version>
</dependency>
` + "```" + `

### Gradle

Add the following dependency to your ` + "`build.gradle`" + `:

` + "```gradle" + `
implementation '%s:%s:%s'
` + "```" + `

## Quick Start

` + "```java" + `
import %s.%s;
import %s.config.ClientConfig;
import %s.models.*;
import %s.exceptions.*;

public class Example {
    public static void main(String[] args) {
        // Initialize the client
        ClientConfig config = new ClientConfig(
            "http://localhost:8080",
            "your-api-key"
        );
        
        %s client = new %s(config);
        
        try {
            // Execute workflow
            Map<String, Object> inputData = new HashMap<>();
            inputData.put("key", "value");
            
            ExecutionResult result = client.executeWorkflow(inputData);
            System.out.println("Execution ID: " + result.getId());
            System.out.println("Status: " + result.getStatus());
            
            // Check execution status
            ExecutionStatus status = client.getExecutionStatus(result.getId());
            System.out.println("Current status: " + status.getStatus());
            System.out.println("Progress: " + status.getProgress() + "%%");
            
        } catch (AuthenticationException e) {
            System.err.println("Invalid API key: " + e.getMessage());
        } catch (ValidationException e) {
            System.err.println("Invalid input: " + e.getMessage());
        } catch (ExecutionException e) {
            System.err.println("Execution failed: " + e.getMessage());
        } catch (MagicFlowException e) {
            System.err.println("Error: " + e.getMessage());
        }
    }
}
` + "```" + `

## API Reference

### Client Methods

#### executeWorkflow

Executes the %s workflow with the provided input.

` + "```java" + `
public ExecutionResult executeWorkflow(Map<String, Object> inputData) 
    throws MagicFlowException
` + "```" + `

#### getExecutionStatus

Retrieves the status of a workflow execution.

` + "```java" + `
public ExecutionStatus getExecutionStatus(UUID executionId) 
    throws MagicFlowException
` + "```" + `

#### cancelExecution

Cancels a running workflow execution.

` + "```java" + `
public void cancelExecution(UUID executionId) 
    throws MagicFlowException
` + "```" + `

#### getExecutionResult

Retrieves the result of a completed workflow execution.

` + "```java" + `
public ExecutionResult getExecutionResult(UUID executionId) 
    throws MagicFlowException
` + "```" + `

#### waitForCompletion

Waits for a workflow execution to complete.

` + "```java" + `
public ExecutionResult waitForCompletion(
    UUID executionId, 
    Duration timeout, 
    Duration pollInterval
) throws MagicFlowException
` + "```" + `

%s

## Models

### ExecutionResult

Represents the result of a workflow execution.

` + "```java" + `
public class ExecutionResult {
    private UUID id;
    private UUID workflowId;
    private ExecutionStatus status;
    private Map<String, Object> input;
    private Map<String, Object> output;
    private String error;
    private Instant startedAt;
    private Instant completedAt;
    private Duration duration;
    
    // Getters and setters...
}
` + "```" + `

### ExecutionStatus

Represents the status of a workflow execution.

` + "```java" + `
public class ExecutionStatus {
    private UUID id;
    private Status status;
    private int progress;
    private String currentStep;
    private List<StepStatus> steps;
    private Instant startedAt;
    private Instant updatedAt;
    
    public enum Status {
        PENDING, RUNNING, COMPLETED, FAILED, CANCELLED
    }
    
    // Getters and setters...
}
` + "```" + `

### StepStatus

Represents the status of a workflow step.

` + "```java" + `
public class StepStatus {
    private String id;
    private String name;
    private Status status;
    private Map<String, Object> input;
    private Map<String, Object> output;
    private String error;
    private Instant startedAt;
    private Instant completedAt;
    private Duration duration;
    
    public enum Status {
        PENDING, RUNNING, COMPLETED, FAILED, SKIPPED
    }
    
    // Getters and setters...
}
` + "```" + `

## Configuration

### ClientConfig

` + "```java" + `
ClientConfig config = new ClientConfig()
    .setBaseUrl("http://localhost:8080")
    .setApiKey("your-api-key")
    .setTimeout(Duration.ofSeconds(30))
    .setRetryAttempts(3)
    .setRetryDelay(Duration.ofSeconds(1))
    .setDebug(false);

%s client = new %s(config);
` + "```" + `

### Environment Variables

` + "```bash" + `
# Set default base URL
export MAGICFLOW_BASE_URL="http://localhost:8080"

# Set default API key
export MAGICFLOW_API_KEY="your-api-key"

# Set default timeout (seconds)
export MAGICFLOW_TIMEOUT="30"

# Enable debug logging
export MAGICFLOW_DEBUG="true"
` + "```" + `

## Exception Handling

The client provides several exception types for different error scenarios:

` + "```java" + `
try {
    ExecutionResult result = client.executeWorkflow(inputData);
} catch (AuthenticationException e) {
    // Invalid API key
    System.err.println("Authentication failed: " + e.getMessage());
} catch (ValidationException e) {
    // Invalid input data
    System.err.println("Validation error: " + e.getMessage());
} catch (ExecutionException e) {
    // Workflow execution failed
    System.err.println("Execution failed: " + e.getMessage());
    System.err.println("Execution ID: " + e.getExecutionId());
    System.err.println("Step ID: " + e.getStepId());
} catch (TimeoutException e) {
    // Request timed out
    System.err.println("Request timed out: " + e.getMessage());
} catch (NetworkException e) {
    // Network connectivity issues
    System.err.println("Network error: " + e.getMessage());
} catch (ApiException e) {
    // General API errors
    System.err.println("API error: " + e.getMessage());
    System.err.println("Status code: " + e.getStatusCode());
} catch (MagicFlowException e) {
    // Base exception for all client errors
    System.err.println("General error: " + e.getMessage());
}
` + "```" + `

## Logging

The client uses SLF4J for logging. Add a logging implementation to your project:

### Logback (recommended)

` + "```xml" + `
<dependency>
    <groupId>ch.qos.logback</groupId>
    <artifactId>logback-classic</artifactId>
    <version>1.4.11</version>
</dependency>
` + "```" + `

### Log4j2

` + "```xml" + `
<dependency>
    <groupId>org.apache.logging.log4j</groupId>
    <artifactId>log4j-slf4j2-impl</artifactId>
    <version>2.20.0</version>
</dependency>
` + "```" + `

## Development

### Building the Project

` + "```bash" + `
# Using Maven
mvn clean compile

# Using Gradle
./gradlew build
` + "```" + `

### Running Tests

` + "```bash" + `
# Using Maven
mvn test

# Using Gradle
./gradlew test
` + "```" + `

### Generating Documentation

` + "```bash" + `
# Using Maven
mvn javadoc:javadoc

# Using Gradle
./gradlew javadoc
` + "```" + `

### Code Coverage

` + "```bash" + `
# Using Maven
mvn jacoco:report

# Using Gradle
./gradlew jacocoTestReport
` + "```" + `

## Requirements

- Java 11 or higher
- Maven 3.6+ or Gradle 7.0+

## License

Generated code - see original workflow license.
`,
		data.Workflow.Name,
		data.Workflow.Name,
		groupId,
		artifactId,
		version,
		groupId,
		artifactId,
		version,
		data.PackageName,
		data.ClassName,
		data.PackageName,
		data.PackageName,
		data.PackageName,
		data.ClassName,
		data.ClassName,
		data.Workflow.Name,
		h.generateMethodDocs(data.Methods),
		data.ClassName,
		data.ClassName,
	)

	return GeneratedFile{
		Path:     "README.md",
		Content:  content,
		Language: "markdown",
		Type:     "documentation",
	}, nil
}

// generateBaseModelContent generates content for base model classes
func (h *JavaHandler) generateBaseModelContent(data *TemplateData, modelName string) string {
	switch modelName {
	case "ExecutionResult":
		return fmt.Sprintf(`package %s.models;

import com.fasterxml.jackson.annotation.JsonProperty;
import java.time.Instant;
import java.time.Duration;
import java.util.Map;
import java.util.UUID;

/**
 * Represents the result of a workflow execution
 * Generated at: %s
 */
public class ExecutionResult {
    @JsonProperty("id")
    private UUID id;
    
    @JsonProperty("workflow_id")
    private UUID workflowId;
    
    @JsonProperty("status")
    private ExecutionStatus.Status status;
    
    @JsonProperty("input")
    private Map<String, Object> input;
    
    @JsonProperty("output")
    private Map<String, Object> output;
    
    @JsonProperty("error")
    private String error;
    
    @JsonProperty("started_at")
    private Instant startedAt;
    
    @JsonProperty("completed_at")
    private Instant completedAt;
    
    @JsonProperty("duration")
    private Long durationMs;
    
    // Constructors
    public ExecutionResult() {}
    
    // Getters and setters
    public UUID getId() { return id; }
    public void setId(UUID id) { this.id = id; }
    
    public UUID getWorkflowId() { return workflowId; }
    public void setWorkflowId(UUID workflowId) { this.workflowId = workflowId; }
    
    public ExecutionStatus.Status getStatus() { return status; }
    public void setStatus(ExecutionStatus.Status status) { this.status = status; }
    
    public Map<String, Object> getInput() { return input; }
    public void setInput(Map<String, Object> input) { this.input = input; }
    
    public Map<String, Object> getOutput() { return output; }
    public void setOutput(Map<String, Object> output) { this.output = output; }
    
    public String getError() { return error; }
    public void setError(String error) { this.error = error; }
    
    public Instant getStartedAt() { return startedAt; }
    public void setStartedAt(Instant startedAt) { this.startedAt = startedAt; }
    
    public Instant getCompletedAt() { return completedAt; }
    public void setCompletedAt(Instant completedAt) { this.completedAt = completedAt; }
    
    public Duration getDuration() {
        return durationMs != null ? Duration.ofMillis(durationMs) : null;
    }
    
    public void setDuration(Duration duration) {
        this.durationMs = duration != null ? duration.toMillis() : null;
    }
}
`,
			data.PackageName,
			data.GeneratedAt.Format("2006-01-02 15:04:05"),
		)

	case "ExecutionStatus":
		return fmt.Sprintf(`package %s.models;

import com.fasterxml.jackson.annotation.JsonProperty;
import java.time.Instant;
import java.util.List;
import java.util.UUID;

/**
 * Represents the status of a workflow execution
 * Generated at: %s
 */
public class ExecutionStatus {
    @JsonProperty("id")
    private UUID id;
    
    @JsonProperty("status")
    private Status status;
    
    @JsonProperty("progress")
    private int progress;
    
    @JsonProperty("current_step")
    private String currentStep;
    
    @JsonProperty("steps")
    private List<StepStatus> steps;
    
    @JsonProperty("started_at")
    private Instant startedAt;
    
    @JsonProperty("updated_at")
    private Instant updatedAt;
    
    public enum Status {
        @JsonProperty("pending") PENDING,
        @JsonProperty("running") RUNNING,
        @JsonProperty("completed") COMPLETED,
        @JsonProperty("failed") FAILED,
        @JsonProperty("cancelled") CANCELLED
    }
    
    // Constructors
    public ExecutionStatus() {}
    
    // Getters and setters
    public UUID getId() { return id; }
    public void setId(UUID id) { this.id = id; }
    
    public Status getStatus() { return status; }
    public void setStatus(Status status) { this.status = status; }
    
    public int getProgress() { return progress; }
    public void setProgress(int progress) { this.progress = progress; }
    
    public String getCurrentStep() { return currentStep; }
    public void setCurrentStep(String currentStep) { this.currentStep = currentStep; }
    
    public List<StepStatus> getSteps() { return steps; }
    public void setSteps(List<StepStatus> steps) { this.steps = steps; }
    
    public Instant getStartedAt() { return startedAt; }
    public void setStartedAt(Instant startedAt) { this.startedAt = startedAt; }
    
    public Instant getUpdatedAt() { return updatedAt; }
    public void setUpdatedAt(Instant updatedAt) { this.updatedAt = updatedAt; }
}
`,
			data.PackageName,
			data.GeneratedAt.Format("2006-01-02 15:04:05"),
		)

	case "StepStatus":
		return fmt.Sprintf(`package %s.models;

import com.fasterxml.jackson.annotation.JsonProperty;
import java.time.Duration;
import java.time.Instant;
import java.util.Map;

/**
 * Represents the status of a workflow step
 * Generated at: %s
 */
public class StepStatus {
    @JsonProperty("id")
    private String id;
    
    @JsonProperty("name")
    private String name;
    
    @JsonProperty("status")
    private Status status;
    
    @JsonProperty("input")
    private Map<String, Object> input;
    
    @JsonProperty("output")
    private Map<String, Object> output;
    
    @JsonProperty("error")
    private String error;
    
    @JsonProperty("started_at")
    private Instant startedAt;
    
    @JsonProperty("completed_at")
    private Instant completedAt;
    
    @JsonProperty("duration")
    private Long durationMs;
    
    public enum Status {
        @JsonProperty("pending") PENDING,
        @JsonProperty("running") RUNNING,
        @JsonProperty("completed") COMPLETED,
        @JsonProperty("failed") FAILED,
        @JsonProperty("skipped") SKIPPED
    }
    
    // Constructors
    public StepStatus() {}
    
    // Getters and setters
    public String getId() { return id; }
    public void setId(String id) { this.id = id; }
    
    public String getName() { return name; }
    public void setName(String name) { this.name = name; }
    
    public Status getStatus() { return status; }
    public void setStatus(Status status) { this.status = status; }
    
    public Map<String, Object> getInput() { return input; }
    public void setInput(Map<String, Object> input) { this.input = input; }
    
    public Map<String, Object> getOutput() { return output; }
    public void setOutput(Map<String, Object> output) { this.output = output; }
    
    public String getError() { return error; }
    public void setError(String error) { this.error = error; }
    
    public Instant getStartedAt() { return startedAt; }
    public void setStartedAt(Instant startedAt) { this.startedAt = startedAt; }
    
    public Instant getCompletedAt() { return completedAt; }
    public void setCompletedAt(Instant completedAt) { this.completedAt = completedAt; }
    
    public Duration getDuration() {
        return durationMs != null ? Duration.ofMillis(durationMs) : null;
    }
    
    public void setDuration(Duration duration) {
        this.durationMs = duration != null ? duration.toMillis() : null;
    }
}
`,
			data.PackageName,
			data.GeneratedAt.Format("2006-01-02 15:04:05"),
		)

	default:
		return fmt.Sprintf(`package %s.models;

/**
 * %s model
 * Generated at: %s
 */
public class %s {
    // TODO: Implement model fields
}
`,
			data.PackageName,
			modelName,
			data.GeneratedAt.Format("2006-01-02 15:04:05"),
			modelName,
		)
	}
}

// generateModelContent generates content for workflow-specific models
func (h *JavaHandler) generateModelContent(data *TemplateData, model ModelData) string {
	content := fmt.Sprintf(`package %s.models;

import com.fasterxml.jackson.annotation.JsonProperty;

/**
 * %s
 * Generated at: %s
 */
public class %s {
`,
		data.PackageName,
		model.Description,
		data.GeneratedAt.Format("2006-01-02 15:04:05"),
		model.Name,
	)

	// Generate fields
	for _, field := range model.Fields {
		javaType := h.mapFieldTypeToJava(field.Type)
		content += fmt.Sprintf(`    @JsonProperty("%s")
    private %s %s;

`,
			ToSnakeCase(field.Name),
			javaType,
			field.Name,
		)
	}

	// Generate constructor
	content += "    // Constructors\n"
	content += fmt.Sprintf("    public %s() {}\n\n", model.Name)

	// Generate getters and setters
	content += "    // Getters and setters\n"
	for _, field := range model.Fields {
		javaType := h.mapFieldTypeToJava(field.Type)
		content += fmt.Sprintf(`    public %s get%s() { return %s; }
    public void set%s(%s %s) { this.%s = %s; }

`,
			javaType,
			ToPascalCase(field.Name),
			field.Name,
			ToPascalCase(field.Name),
			javaType,
			field.Name,
			field.Name,
			field.Name,
		)
	}

	content += "}"
	return content
}

// generateExceptionContent generates content for exception classes
func (h *JavaHandler) generateExceptionContent(data *TemplateData, exceptionName string) string {
	switch exceptionName {
	case "MagicFlowException":
		return fmt.Sprintf(`package %s.exceptions;

/**
 * Base exception for Magic Flow client
 * Generated at: %s
 */
public class MagicFlowException extends Exception {
    public MagicFlowException(String message) {
        super(message);
    }
    
    public MagicFlowException(String message, Throwable cause) {
        super(message, cause);
    }
}
`,
			data.PackageName,
			data.GeneratedAt.Format("2006-01-02 15:04:05"),
		)

	case "ApiException":
		return fmt.Sprintf(`package %s.exceptions;

import java.util.Map;

/**
 * Exception for API errors
 * Generated at: %s
 */
public class ApiException extends MagicFlowException {
    private final int statusCode;
    private final Map<String, Object> responseData;
    
    public ApiException(String message, int statusCode) {
        super(message);
        this.statusCode = statusCode;
        this.responseData = null;
    }
    
    public ApiException(String message, int statusCode, Map<String, Object> responseData) {
        super(message);
        this.statusCode = statusCode;
        this.responseData = responseData;
    }
    
    public int getStatusCode() {
        return statusCode;
    }
    
    public Map<String, Object> getResponseData() {
        return responseData;
    }
}
`,
			data.PackageName,
			data.GeneratedAt.Format("2006-01-02 15:04:05"),
		)

	case "ExecutionException":
		return fmt.Sprintf(`package %s.exceptions;

import java.util.UUID;

/**
 * Exception for workflow execution errors
 * Generated at: %s
 */
public class ExecutionException extends MagicFlowException {
    private final UUID executionId;
    private final String stepId;
    
    public ExecutionException(String message) {
        super(message);
        this.executionId = null;
        this.stepId = null;
    }
    
    public ExecutionException(String message, UUID executionId) {
        super(message);
        this.executionId = executionId;
        this.stepId = null;
    }
    
    public ExecutionException(String message, UUID executionId, String stepId) {
        super(message);
        this.executionId = executionId;
        this.stepId = stepId;
    }
    
    public UUID getExecutionId() {
        return executionId;
    }
    
    public String getStepId() {
        return stepId;
    }
}
`,
			data.PackageName,
			data.GeneratedAt.Format("2006-01-02 15:04:05"),
		)

	default:
		return fmt.Sprintf(`package %s.exceptions;

/**
 * %s
 * Generated at: %s
 */
public class %s extends MagicFlowException {
    public %s(String message) {
        super(message);
    }
    
    public %s(String message, Throwable cause) {
        super(message, cause);
    }
}
`,
			data.PackageName,
			exceptionName,
			data.GeneratedAt.Format("2006-01-02 15:04:05"),
			exceptionName,
			exceptionName,
			exceptionName,
		)
	}
}

// generateImports generates the list of imports needed
func (h *JavaHandler) generateImports(workflow *models.Workflow, request *GenerationRequest) []string {
	imports := []string{
		"java.util.Map",
		"java.util.HashMap",
		"java.util.List",
		"java.util.UUID",
		"java.time.Instant",
		"java.time.Duration",
		"com.fasterxml.jackson.databind.ObjectMapper",
		"okhttp3.OkHttpClient",
		"okhttp3.Request",
		"okhttp3.Response",
	}

	// Add test imports if tests are included
	if request.IncludeTests {
		imports = append(imports, "org.junit.jupiter.api.Test", "org.mockito.Mock", "org.mockito.MockitoAnnotations")
	}

	return imports
}

// generateModels generates model definitions from workflow
func (h *JavaHandler) generateModels(workflow *models.Workflow) []ModelData {
	var models []ModelData

	// Generate models based on workflow inputs/outputs
	if workflow.Definition.Input != nil {
		for key, schema := range workflow.Definition.Input {
			model := ModelData{
				Name:        ToPascalCase(key) + "Input",
				Description: fmt.Sprintf("Input model for %s", key),
				Fields:      h.generateFieldsFromSchema(schema),
			}
			models = append(models, model)
		}
	}

	if workflow.Definition.Output != nil {
		for key, schema := range workflow.Definition.Output {
			model := ModelData{
				Name:        ToPascalCase(key) + "Output",
				Description: fmt.Sprintf("Output model for %s", key),
				Fields:      h.generateFieldsFromSchema(schema),
			}
			models = append(models, model)
		}
	}

	return models
}

// generateFieldsFromSchema generates field definitions from schema
func (h *JavaHandler) generateFieldsFromSchema(schema interface{}) []FieldData {
	var fields []FieldData

	// This is a simplified implementation
	// In a real scenario, you would parse the JSON schema properly
	if schemaMap, ok := schema.(map[string]interface{}); ok {
		if properties, ok := schemaMap["properties"].(map[string]interface{}); ok {
			for fieldName, fieldSchema := range properties {
				field := FieldData{
					Name:        ToPascalCase(fieldName),
					Type:        h.mapSchemaTypeToJavaType(fieldSchema),
					Description: h.getSchemaDescription(fieldSchema),
					Required:    h.isFieldRequired(fieldName, schemaMap),
				}
				fields = append(fields, field)
			}
		}
	}

	return fields
}

// mapSchemaTypeToJavaType maps JSON schema types to Java types
func (h *JavaHandler) mapSchemaTypeToJavaType(schema interface{}) string {
	if schemaMap, ok := schema.(map[string]interface{}); ok {
		if schemaType, ok := schemaMap["type"].(string); ok {
			switch schemaType {
			case "string":
				return "String"
			case "integer":
				return "Integer"
			case "number":
				return "Double"
			case "boolean":
				return "Boolean"
			case "array":
				return "List<Object>"
			case "object":
				return "Map<String, Object>"
			}
		}
	}
	return "Object"
}

// mapFieldTypeToJava maps field types to Java types
func (h *JavaHandler) mapFieldTypeToJava(fieldType string) string {
	switch fieldType {
	case "string":
		return "String"
	case "integer", "int":
		return "Integer"
	case "number", "float", "double":
		return "Double"
	case "boolean", "bool":
		return "Boolean"
	case "array", "list":
		return "List<Object>"
	case "object", "map":
		return "Map<String, Object>"
	case "uuid":
		return "UUID"
	case "datetime", "timestamp":
		return "Instant"
	case "duration":
		return "Duration"
	default:
		return "Object"
	}
}

// getSchemaDescription extracts description from schema
func (h *JavaHandler) getSchemaDescription(schema interface{}) string {
	if schemaMap, ok := schema.(map[string]interface{}); ok {
		if desc, ok := schemaMap["description"].(string); ok {
			return desc
		}
	}
	return ""
}

// isFieldRequired checks if a field is required
func (h *JavaHandler) isFieldRequired(fieldName string, schema map[string]interface{}) bool {
	if required, ok := schema["required"].([]interface{}); ok {
		for _, req := range required {
			if reqStr, ok := req.(string); ok && reqStr == fieldName {
				return true
			}
		}
	}
	return false
}

// generateMethodDocs generates documentation for methods
func (h *JavaHandler) generateMethodDocs(methods []MethodData) string {
	if len(methods) == 0 {
		return ""
	}

	docs := "\n### Workflow Methods\n\n"
	for _, method := range methods {
		docs += fmt.Sprintf("#### %s\n\n%s\n\n```java\npublic ExecutionResult %s(Map<String, Object> inputData) \n    throws MagicFlowException\n```\n\n",
			method.Name,
			method.Description,
			ToCamelCase(method.Name),
		)
	}

	return docs
}

// isValidJavaPackageName validates Java package name
func isValidJavaPackageName(packageName string) bool {
	if packageName == "" {
		return false
	}

	// Split by dots and validate each part
	parts := strings.Split(packageName, ".")
	for _, part := range parts {
		if !isValidJavaIdentifier(part) {
			return false
		}
	}

	return true
}

// isValidJavaIdentifier validates Java identifier
func isValidJavaIdentifier(identifier string) bool {
	if identifier == "" {
		return false
	}

	// Check if it's a Java keyword
	javaKeywords := map[string]bool{
		"abstract": true, "assert": true, "boolean": true, "break": true,
		"byte": true, "case": true, "catch": true, "char": true,
		"class": true, "const": true, "continue": true, "default": true,
		"do": true, "double": true, "else": true, "enum": true,
		"extends": true, "final": true, "finally": true, "float": true,
		"for": true, "goto": true, "if": true, "implements": true,
		"import": true, "instanceof": true, "int": true, "interface": true,
		"long": true, "native": true, "new": true, "package": true,
		"private": true, "protected": true, "public": true, "return": true,
		"short": true, "static": true, "strictfp": true, "super": true,
		"switch": true, "synchronized": true, "this": true, "throw": true,
		"throws": true, "transient": true, "try": true, "void": true,
		"volatile": true, "while": true,
	}

	if javaKeywords[identifier] {
		return false
	}

	// Check first character
	firstChar := rune(identifier[0])
	if !((firstChar >= 'a' && firstChar <= 'z') || (firstChar >= 'A' && firstChar <= 'Z') || firstChar == '_' || firstChar == '$') {
		return false
	}

	// Check remaining characters
	for _, char := range identifier[1:] {
		if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') || char == '_' || char == '$') {
			return false
		}
	}

	return true
}