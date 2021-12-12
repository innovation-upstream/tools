# gen code frame

The goal of this project is to define a metaprogramming spec and generic code
generator that is capable of taking abstract model definitions and
relationships, and using them to generate working implementation in any
language.  Similar to how Kubernetes defines a spec and "container generator"
for container creation and runtime environments.

The high-level components of this repo is as follows:

- Models
  - Core Model Functions
  - Hooks
- Generic Implementation Modules
- Specialized Hook Effects
- Code Generator

## Model

A model describes a data model of any kind (ex. a User) in terms of its
functionality (ex. how it can be queried).  All of a models functionality is
defined by hooks into a "Core Model Function".


## Core Model Functions

Core model functions are a set of simple functions that all models have.
The modules used to generate implementation for a model is responsible for the
exact behaviour of each core model function.  This projects defines the
following core model functions:

**Query**

A function that searchs the persistance layer for a set of models based on their
attributes and relationships.

**Mutate**

A change to a model in the persistance layer.

**Translate**

A function that takes the model as input, and returns some derivitive of the
model as the output.  Translations can be lossy. (ex. Format an appointment to
be Human-readable)


## Hooks

Hooks are called during the execution of a core model function.  Hooks have
varying capabilities depending on the core model function.  Hooks are
implemented in the target language, by the end developer (person using this
tool).  Hooks are inserted into the generated code by the modules used to
generate it.

**PreMutate**

Called before a model is mutated, this could be a create or an update. Can only
impact control flow by returning an error, which cancels execution.

**PostMutate**

Called after a model is mutated, this could be a create or an update. Can only
impact control flow by returning an error, which cancels execution.

**PreQuery**

Called before a model is queried. Can only impact control flow by returning an
error, which cancels execution.

**PostQuery**

Called after a model is queried. Can only impact control flow by returning an
error, which cancels execution.
