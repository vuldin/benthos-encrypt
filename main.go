package main

import (
  "context"

  "github.com/benthosdev/benthos/v4/public/service"

  // Import only pure and standard io Benthos components
  // TODO Why only io and pure? How do I know which components to import?
  _ "github.com/benthosdev/benthos/v4/public/components/io"
  _ "github.com/benthosdev/benthos/v4/public/components/pure"

  // In order to import _all_ Benthos components for third party services
  // uncomment the following line:
  // _ "github.com/benthosdev/benthos/v4/public/components/all"

  //_ "github.com/vuldin/benthos-encrypt/processor"
)

func main() {
  service.RunCLI(context.Background())
}
