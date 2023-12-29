package db

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type TestCryptoModel struct {
	FieldA string  `gorm:"column:field_1" encrypt:"true"`
	FieldB *string `encrypt:"aes"`
	FieldC string
	FieldD *string
	FieldE []byte `encrypt:"aes"`
}

func TestCryptoHook(t *testing.T) {
	mashal := func(ctx context.Context, tag string, val string) (string, error) {
		if val == "" {
			return val, nil
		}
		return val + "_encrypted", nil
	}
	unmashal := func(ctx context.Context, tag string, val string) (string, error) {
		if val == "" {
			return val, nil
		}
		return val + "_decrypted", nil
	}
	dial := WithInitializeHook(testdb_dial(t), CryptoHook(mashal, unmashal))
	p := testdb_newprovider_with_dial(t, dial, "crypto_hook")
	ctx := context.Background()

	if err := p.UseDB(ctx).AutoMigrate(&TestCryptoModel{}); err != nil {
		t.Fatal(err)
	}
	strptr := func(s string) *string { return &s }

	cases := []struct {
		give *TestCryptoModel
		exp  *TestCryptoModel
	}{
		{
			give: &TestCryptoModel{
				FieldA: "value_a",
				FieldB: strptr("value_b"),
				FieldC: "value_c",
				FieldD: strptr("value_d"),
				FieldE: []byte("value_e"),
			},
			exp: &TestCryptoModel{
				FieldA: "value_a_encrypted_decrypted",
				FieldB: strptr("value_b_encrypted_decrypted"),
				FieldC: "value_c",
				FieldD: strptr("value_d"),
				FieldE: []byte("value_e_encrypted_decrypted"),
			},
		},
		{
			give: &TestCryptoModel{},
			exp:  &TestCryptoModel{},
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("case %d", i+1), func(t *testing.T) {
			if err := p.UseDB(ctx).Create(c.give).Error; err != nil {
				t.Error(err)
			}
			if !cmp.Equal(c.give, c.exp) {
				t.Errorf("diff: %s", cmp.Diff(c.exp, c.give))
			}
		})
	}
}

func TestCryptoHook_Mashal(t *testing.T) {
	mashal := func(ctx context.Context, tag string, val string) (string, error) {
		if val == "" {
			return val, nil
		}
		return val + "_encrypted", nil
	}
	unmashal := func(ctx context.Context, tag string, val string) (string, error) {
		return val, nil
	}
	dial := WithInitializeHook(testdb_dial(t), CryptoHook(mashal, unmashal))
	p := testdb_newprovider_with_dial(t, dial, "crypto_hook")
	ctx := context.Background()

	if err := p.UseDB(ctx).AutoMigrate(&TestCryptoModel{}); err != nil {
		t.Fatal(err)
	}
	strptr := func(s string) *string { return &s }

	cases := []struct {
		give *TestCryptoModel
		exp  *TestCryptoModel
	}{
		{
			give: &TestCryptoModel{
				FieldA: "value_a",
				FieldB: strptr("value_b"),
				FieldC: "value_c",
				FieldD: strptr("value_d"),
				FieldE: []byte("value_e"),
			},
			exp: &TestCryptoModel{
				FieldA: "value_a_encrypted",
				FieldB: strptr("value_b_encrypted"),
				FieldC: "value_c",
				FieldD: strptr("value_d"),
				FieldE: []byte("value_e_encrypted"),
			},
		},
		{
			give: &TestCryptoModel{},
			exp:  &TestCryptoModel{},
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("case %d", i+1), func(t *testing.T) {
			if err := p.UseDB(ctx).Create(c.give).Error; err != nil {
				t.Error(err)
			}
			if !cmp.Equal(c.give, c.exp) {
				t.Errorf("diff: %s", cmp.Diff(c.exp, c.give))
			}
		})
	}
}

func TestEncryptAndDecrypt(t *testing.T) {
	aesKey := func(ctx context.Context) (string, error) {
		return strings.Repeat("a", 16), nil
	}
	encrypt := EncryptTagHandler(aesKey)
	decrypt := DecryptTagHandler(aesKey)

	ctx := context.Background()
	raw := "raw data"
	en, err := encrypt(ctx, "aes", raw)
	if err != nil {
		t.Fatal(err)
	}
	if raw == en {
		t.Errorf("expect not equal %s , %s", raw, en)
	}
	de, err := decrypt(ctx, "aes", en)
	if err != nil {
		t.Fatal(err)
	}
	if de != raw {
		t.Errorf("expect: %s, got: %s", raw, de)
	}
}
