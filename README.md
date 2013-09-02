# Cloudburst


Cloudburst is an efficient workload generation toolkit for cloud environments developed in Go programming language. It uses and applies workload distributions to generate load on servers. Different distributions can be configured to aim at different target servers. For the actual load generation, different operations can be specified to provide and compose a realistic behavior.

## Architecture

![Cloudburst Architecture](https://github.com/johanneskross/cloudburst/blob/refactoring/cloudburst_architecture.png?raw=true "Cloudburst Architecture")

* A Benchmark is the starting point of the system. 

* A Scenario represents a concrete experiment and coordinates it.

* A Target Schedule specifies the experiment settings and contains one or more Target Configurations.

* A Target Configuration represents a configuration for one target.

* A Target Manager uses the Target Schedule and a Target Factory to create Targets.

* A Target exists for each target server on which workload will be created. It uses a Load Manager to schedule the amount and duration of Agents.

* A Load Manager maintains a workload distribution.

* An Agent embodies an user and executes operations on an target server. It receives a Generator implementation of its parent Target which comprises a method to request an Operation. An Operation again comprises a run-method that is then executed by the Agent. 

* A Scoreboard receives statistics after each operation execution from agents of one target.

## Implementation

A sample implementation can be found at https://github.com/johanneskross/benchmark/
