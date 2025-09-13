# Integration

Integration module includes integration tests for `github.com/hossein1376/spallet`.
It starts a database container, run the tests, and clean it up. Also, end-to-end
(E2E) tests can be added here as well.

Make sure you have Docker or Podman installed beforehand and it's running. Also,
keep in mind that this package was written in a Linux machine and may or may not
work in other operating services.

## Reason behind it

Main reason to have a separate module for such tests is to keep the project's
dependency list as concise as possible.  With this approach, the functionality
of packages and functions can be tested with a real database or from a real
client's point of view.

Not being able to import from the main module's internal folder or not having
any access to private functions and fields might prove to be a challenge on its
own, but keep in mind that the intention behind these tests is to make sure the
application is acting as expected. Test behaviour, not the implementation, so to
speak.
