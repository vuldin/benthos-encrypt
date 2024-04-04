package processor

import (
  "context"
  "crypto/aes"
  "crypto/cipher"
  "crypto/rand"
  "encoding/base64"
  "encoding/hex"
  "io"

  "github.com/benthosdev/benthos/v4/public/service"
  "github.com/tidwall/gjson"
  "github.com/tidwall/sjson"
)

func init() {
  configSpec := service.NewConfigSpec().
    Version("0.0.1").
    Summary("Encrypts specific fields with a provided key").
    Field(service.NewStringListField("fields").
      Description("List of fields to encrypt.")).
    Field(service.NewStringField("keyString").
      Description("The key used to encrypt."))

  constructor := func(conf *service.ParsedConfig, mgr *service.Resources) (service.Processor, error) {
    return newEncryptProcessor(conf, mgr.Logger()), nil
  }

  err := service.RegisterProcessor("encrypt", configSpec, constructor)
  if err != nil {
    panic(err)
  }
}

//------------------------------------------------------------------------------

type encryptProcessor struct {
  conf   *service.ParsedConfig
  logger *service.Logger
}

func newEncryptProcessor(conf *service.ParsedConfig, logger *service.Logger) *encryptProcessor {
  return &encryptProcessor{
    conf:   conf,
    logger: logger,
  }
}

func (r *encryptProcessor) Process(ctx context.Context, m *service.Message) (service.MessageBatch, error) {
  bytesContent, err := m.AsBytes()
  if err != nil {
    return nil, err
  }

  keyString, err := r.conf.FieldString("keyString")
  fields, err := r.conf.FieldStringList("fields")
  for _, field := range fields {
    textBytes := gjson.GetBytes(bytesContent, field)
    text := textBytes.String()
    cryptoText := encrypt(keyString, text)
    value, _ := sjson.SetBytes(bytesContent, field, cryptoText)
    bytesContent = value
  }
  m.SetBytes(bytesContent)
  return []*service.Message{m}, nil
}

func (r *encryptProcessor) Close(ctx context.Context) error {
  return nil
}

func encrypt(keyString string, stringToEncrypt string) (encryptedString string) {
  key, _ := hex.DecodeString(keyString)
  plaintext := []byte(stringToEncrypt)

  block, err := aes.NewCipher(key)
  if err != nil {
    panic(err.Error())
  }

  ciphertext := make([]byte, aes.BlockSize+len(plaintext))
  iv := ciphertext[:aes.BlockSize]
  if _, err := io.ReadFull(rand.Reader, iv); err != nil {
    panic(err)
  }

  stream := cipher.NewCFBEncrypter(block, iv)
  stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

  return base64.URLEncoding.EncodeToString(ciphertext)
}

