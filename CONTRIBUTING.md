# Contribution guide

This contribution guide will focus mostly on creating a new backend for the DAMAP integration. It assumes that you have at least a passing familiarity with Go.

## Creating a new backend

In order to create a backend, you must implement the interfaces described in the [`backend`](backend/backend.go) package. First, you must implement the `Descriptor` interface. Its main purpose is to provide data about the backend you are creating. Its most important function is to provide the configuration struct that the integration server can fill with options from the CLI or environment using the `GetConfiguration` function.

Once you have created the implementation for this interface take a moment to do two things:

1. Create an `init()` function, which will be called automatically when your backend package is first imported.
2. Add your backend package to [`backends.go`](backends.go) as an import.

Now that this is done, create your configuration struct. You can and should use several tags here:

```go
type configuration struct {
    SomeOption string `default:"some value" cli:"some-option" env:"SOME_OPTION" help:"Some help text here"`
}
```

These tags will let the DAMAP integration automatically populate your configuration struct from the CLI or environment variables. Now there is only one thing left to do: you need to add a `Configure()` receiver function to this configuration struct and then implement the `ConfiguredBackend` interface that actually fetches your data.

If you did everything correctly, your backend should now show up in the supported backends list along with its options.