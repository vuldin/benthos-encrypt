package processor

import (
  "bytes"
  "context"

  "github.com/benthosdev/benthos/v4/public/service"
)

func init() {
  // Config spec is empty for now as we don't have any dynamic fields.
  configSpec := service.NewConfigSpec()
    .Version("0.0.1")
    .Summary("Encrypts specific fields with a provided key")
    .Field(service.NewStringListField("fields")
      .Description("List of fields to encrypt."))
    .Field(service.NewStringField("keyString")
      .Description("The key used to encrypt."))

  constructor := func(conf *service.ParsedConfig, mgr *service.Resources) (service.Processor, error) {
    return newReverseProcessor(mgr.Logger()), nil
  }

  err := service.RegisterProcessor("encrypt", configSpec, constructor)
  if err != nil {
    panic(err)
  }
}

//------------------------------------------------------------------------------

type encryptProcessor struct {
  logger           *service.Logger
}

func newEncryptProcessor(logger *service.Logger) *encryptProcessor {
  return &reverseProcessor{
    logger:           logger,
  }
}

func (r *encryptProcessor) Process(ctx context.Context, m *service.Message) (service.MessageBatch, error) {
//  bytesContent, err := m.AsBytes()
//  if err != nil {
//    return nil, err
//  }

//  newBytes := make([]byte, len(bytesContent))
//  for i, b := range bytesContent {
//    newBytes[len(newBytes)-i-1] = b
//  }
//
//  if bytes.Equal(newBytes, bytesContent) {
//    r.logger.Infof("Woah! This is like totally a palindrome: %s", bytesContent)
//    r.countPalindromes.Incr(1)
//  }

//  m.SetBytes(newBytes)
  return []*service.Message{m}, nil
}

func (r *encryptProcessor) Close(ctx context.Context) error {
  return nil
}

