package gw

import "testing"

func TestDefaultPasswordProtectMd5Impl_Sign(t *testing.T) {
	var def = DefaultPasswordSignerMd5Impl{
		salt: "cvMC33eY7o9YKarcUr7VCf9XLFmHXKWJ",
	}
	t.Logf(def.Sign("123@456"))
}
