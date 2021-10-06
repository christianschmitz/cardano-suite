package network

import (
  "errors"
)

type HandshakeMessage interface { // all messages of a given protocol need to have the same interface type
}

type HandshakeProposeVersions struct {
  Versions map[int]interface{}
}

type HandshakeAcceptVersion struct {
  Version int
  Params   interface{}
}

type HandshakeError interface {
  HandshakeMessage
  Error() error
}

type HandshakeVersionMismatch struct {
  ValidVersions []int
}

type HandshakeDecodeError struct {
  Version int
  Reason  string
}

type HandshakeRefused struct {
  Version int
  Reason  string
}

// NOTE: the primitive type-names below are different, because the builtin cast functions are not very robust
//    see ../common/cbor_primitive_types.go for the possible cast functions

//go:generate ../cbor-type HandshakeMessage *HandshakeProposeVersions *HandshakeAcceptVersion HandshakeError
//go:generate ../cbor-type *HandshakeProposeVersions "Versions common.IntToInterfMap"
//go:generate ../cbor-type *HandshakeAcceptVersion   "Version common.Int, Params common.Interf"
//go:generate ../cbor-type HandshakeError *HandshakeVersionMismatch *HandshakeDecodeError *HandshakeRefused
//go:generate ../cbor-type *HandshakeVersionMismatch "ValidVersions common.IntList"
//go:generate ../cbor-type *HandshakeDecodeError     "Version common.Int, Reason common.String"
//go:generate ../cbor-type *HandshakeRefused         "Version common.Int, Reason common.String"

func NewHandshakeProposeVersions(magic int) *HandshakeProposeVersions {
  return &HandshakeProposeVersions{DefaultProtocolVersions(magic)}
}

func (m *HandshakeProposeVersions) Intersect(version map[int]interface{}) (int, interface{}, error) {
  vIntersect := -1
  var pIntersect interface{} = nil

  for v, p := range m.Versions {
    if _, ok := version[v]; ok {
      if v > vIntersect {
        vIntersect = v
        pIntersect = p
      }
    }
  }

  if vIntersect == -1 {
    return 0, nil, errors.New("no version intersection found")
  }

  return vIntersect, pIntersect, nil
}

func (m *HandshakeVersionMismatch) Error() error {
  // TODO: format the alternative possible versions
  return errors.New("Handshake version mismatch")
}

func (m *HandshakeDecodeError) Error() error {
  return errors.New("HandshakeDecodeError: " + m.Reason)
}

func (m *HandshakeRefused) Error() error {
  return errors.New("HandshakeRefused: " + m.Reason)
}

func DefaultProtocolVersions(magic int) map[int]interface{} {
  v := make(map[int]interface{})

  v[1] = magic
  v[2] = magic
  v[3] = magic

  l4 := make([]interface{}, 2)
  l4[0] = magic
  l4[1] = false
  v[4] = l4

  l5 := make([]interface{}, 2)
  l5[0] = magic
  l5[1] = false
  v[5] = l5

  l6 := make([]interface{}, 2)
  l6[0] = magic
  l6[1] = false
  v[6] = l6

  return v
}