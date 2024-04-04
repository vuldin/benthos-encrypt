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
  //r.logger.Infof("keyString %s", keyString)
  fields, err := r.conf.FieldStringList("fields")
  //r.logger.Infof("first field %s", fields[0])
  //r.logger.Infof("%s", bytesContent)
  //name := gjson.Get(string(bytesContent), fields[0])
  name := gjson.GetBytes(bytesContent, fields[0])
  //println(value.String())
  //r.logger.Infof(name.String())
  cryptoText := encrypt(keyString, name.String())
  //r.logger.Infof(cryptoText)
  value, _ := sjson.SetBytes(bytesContent, fields[0], cryptoText)
  //r.logger.Infof("%s", value)
  m.SetBytes(value)
  return []*service.Message{m}, nil
}

func (r *encryptProcessor) Close(ctx context.Context) error {
  return nil
}

func encrypt(keyString string, stringToEncrypt string) (encryptedString string) {
  // convert key to bytes
  key, _ := hex.DecodeString(keyString)
  plaintext := []byte(stringToEncrypt)

  //Create a new Cipher Block from the key
  block, err := aes.NewCipher(key)
  if err != nil {
    panic(err.Error())
  }

  // The IV needs to be unique, but not secure. Therefore it's common to
  // include it at the beginning of the ciphertext.
  ciphertext := make([]byte, aes.BlockSize+len(plaintext))
  iv := ciphertext[:aes.BlockSize]
  if _, err := io.ReadFull(rand.Reader, iv); err != nil {
    panic(err)
  }

  stream := cipher.NewCFBEncrypter(block, iv)
  stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

  // convert to base64
  return base64.URLEncoding.EncodeToString(ciphertext)
}

