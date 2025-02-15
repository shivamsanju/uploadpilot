# Temporal DSL Workflow Guide

## Introduction

Temporal is an open-source workflow orchestration engine that allows developers to build scalable and reliable workflows. Temporal DSL (Domain-Specific Language) provides a declarative way to define workflows using YAML syntax, making it easier to manage complex workflows without writing extensive code.

This guide explains **everything from scratch**, ensuring that even those new to Temporal can understand and build workflows using DSL syntax.

---

## 1. Key Concepts

Before diving into the DSL syntax, let's cover the fundamental concepts of Temporal workflows:

### a) Workflow

A **workflow** is a sequence of steps (or activities) executed in a defined order. Workflows can have **sequential** and **parallel** execution.

### b) Activities

**Activities** are the building blocks of a workflow. They are individual tasks that take input parameters and return output results.

### c) Execution Flow

Temporal DSL allows defining workflows in a structured YAML format, ensuring execution follows a predefined sequence.

### d) Parallel Execution

Some workflows require multiple activities to run in parallel. Temporal DSL provides a way to define parallel execution branches.

### e) Variables

Variables store reusable values that can be referenced throughout the workflow.

---

## 2. Temporal DSL Syntax Overview

Temporal DSL uses a **YAML-based** structure to define workflows. Below are the core components:

### a) Variables

```yaml
variables:
  input1: value1
  input2: value2
```

- Defines global variables that can be used throughout the workflow.

### b) Workflow Structure

```yaml
root:
  sequence:
    elements:
```

- `root`: The entry point of the workflow.
- `sequence`: Defines a **step-by-step execution** of activities.
- `elements`: Contains the individual steps to execute.

### c) Activity Definition

```yaml
- activity:
    name: ActivityName
    arguments:
      - input1
    result: output1
```

- `activity`: Defines a **single executable task**.
- `name`: The identifier for the activity.
- `arguments`: List of input parameters.
- `result`: Stores the output of the activity.

### d) Parallel Execution

```yaml
- parallel:
    branches:
      - sequence:
          elements:
            - activity:
                name: ActivityA
                arguments:
                  - input1
                result: outputA
      - sequence:
          elements:
            - activity:
                name: ActivityB
                arguments:
                  - input2
                result: outputB
```

- `parallel`: Defines **multiple execution paths** that run at the same time.
- `branches`: Contains **independent sequences** of activities.

### e) Conditional Execution (if-else logic)

```yaml
- condition:
    if: "${outputA == 'success'}"
    then:
      - activity:
          name: SuccessActivity
          arguments:
            - outputA
    else:
      - activity:
          name: FailureActivity
          arguments:
            - outputA
```

- `condition`: Executes a specific block based on a condition.
- `if`: Specifies the condition to evaluate.
- `then`: Defines the steps to execute if the condition is **true**.
- `else`: Defines steps for the **false** case.

### f) Loop Execution (Iterating Over Items)

```yaml
- loop:
    items: "${inputList}"
    as: item
    do:
      - activity:
          name: ProcessItem
          arguments:
            - "${item}"
```

- `loop`: Iterates over a list of values.
- `items`: The list of values to loop through.
- `as`: The variable representing the current item.
- `do`: Defines the activities to execute for each item.

---

## 3. Full Example Workflow

Here’s a complete workflow demonstrating **sequential execution, parallel execution, and conditions**:

```yaml
variables:
  inputA: "valueA"
  inputB: "valueB"

root:
  sequence:
    elements:
      - activity:
          name: StartTask
          arguments:
            - inputA
          result: resultA

      - parallel:
          branches:
            - sequence:
                elements:
                  - activity:
                      name: Task1
                      arguments:
                        - resultA
                      result: result1
            - sequence:
                elements:
                  - activity:
                      name: Task2
                      arguments:
                        - resultA
                      result: result2

      - condition:
          if: "${result1 == 'success'}"
          then:
            - activity:
                name: SuccessHandler
                arguments:
                  - result1
          else:
            - activity:
                name: FailureHandler
                arguments:
                  - result1

      - activity:
          name: FinalTask
          arguments:
            - result2
          result: finalResult
```

### Execution Flow

1. `StartTask` runs first.
2. Two parallel sequences execute `Task1` and `Task2`.
3. A condition checks if `Task1` succeeded, executing either `SuccessHandler` or `FailureHandler`.
4. `FinalTask` executes using the result of `Task2`.

---

## 4. Best Practices

### ✅ Use Descriptive Activity Names

Always name your activities clearly, so the workflow is easy to understand.

### ✅ Organize Workflow into Logical Steps

Break down workflows into **manageable sequences** instead of one large block.

### ✅ Use Parallel Execution Efficiently

Use parallel blocks for **independent** tasks to optimize execution time.

### ✅ Validate Input & Output Variables

Make sure each activity receives the **correct inputs** and produces **expected outputs**.

### ✅ Handle Conditional Logic Properly

Use **if-else conditions** to handle different execution paths.

---

## 5. Debugging & Troubleshooting

### 🔹 Common Errors

| Error Type             | Cause                       | Solution                              |
| ---------------------- | --------------------------- | ------------------------------------- |
| Missing Variables      | Undefined input variable    | Ensure all referenced variables exist |
| Syntax Error           | YAML formatting issue       | Validate YAML structure               |
| Invalid Execution Flow | Results not properly linked | Check input-output dependencies       |

### 🔹 Debugging Tips

1. **Print Execution Logs** - Use logging to track execution steps.
2. **Test with Simple Workflows** - Start with small workflows before building complex ones.
3. **Validate DSL Syntax** - Use a YAML validator to check for errors.

---

## Conclusion

Temporal DSL provides a **powerful and structured** way to define workflows using YAML. By understanding **activities, sequences, parallel execution, conditions, and loops**, you can create complex and reliable workflows with ease.

Use this guide as a reference when building your Temporal workflows. Happy coding! 🚀
