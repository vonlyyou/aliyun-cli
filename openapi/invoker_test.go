package openapi

import (
	"bufio"
	"github.com/aliyun/aliyun-cli/cli"
	"github.com/aliyun/aliyun-cli/config"
	"github.com/aliyun/aliyun-cli/meta"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBasicInvoker_Init(t *testing.T) {
	cp := &config.Profile{
		Mode: config.AuthenticateMode("AK"),
	}
	invooker := NewBasicInvoker(cp)
	client := invooker.getClient()
	assert.Nil(t, client)

	req := invooker.getRequest()
	assert.Nil(t, req)

	w := new(bufio.Writer)
	ctx := cli.NewCommandContext(w)

	regionflag := config.NewRegionFlag()
	regionflag.SetAssigned(true)
	regionflag.SetValue("cn-hangzhou")
	ctx.Flags().Add(regionflag)

	endpointflag := NewEndpointFlag()
	endpointflag.SetAssigned(true)
	endpointflag.SetValue("ecs.cn-hangzhou.aliyuncs")
	ctx.Flags().Add(endpointflag)

	versionflag := NewVersionFlag()
	versionflag.SetAssigned(true)
	versionflag.SetValue("v1.0")
	ctx.Flags().Add(versionflag)

	headerflag := NewHeaderFlag()
	headerflag.SetValues([]string{"Accept=xml", "Accept=json", "Content-Type=json", "testfail"})
	ctx.Flags().Add(headerflag)

	ctx.Flags().Add(config.NewSkipSecureVerify())

	product := &meta.Product{}
	invooker.profile.Mode = config.AuthenticateMode("DEFAULT")
	err := invooker.Init(ctx, product)
	assert.NotNil(t, err)
	assert.Equal(t, "init client failed unexcepted certificate mode: DEFAULT", err.Error())

	invooker.profile.Mode = config.AuthenticateMode("StsToken")
	err = invooker.Init(ctx, product)
	assert.NotNil(t, err)
	assert.Equal(t, "invaild flag --header `testfail` use `--header HeaderName=Value`", err.Error())

	regionflag.SetAssigned(false)
	endpointflag.SetAssigned(false)
	versionflag.SetAssigned(false)
	headerflag.SetValues([]string{})
	err = invooker.Init(ctx, product)
	assert.NotNil(t, err)
	assert.Equal(t, "missing version for product ", err.Error())

	invooker.profile.Mode = config.AuthenticateMode("StsToken")
	product.Version = "v1.0"
	err = invooker.Init(ctx, product)
	assert.NotNil(t, err)
	assert.Equal(t, "missing region for product ", err.Error())

	invooker.profile.RegionId = "cn-hangzhou"
	err = invooker.Init(ctx, product)
	assert.NotNil(t, err)
	assert.Equal(t, "unknown endpoint for /cn-hangzhou! failed unknown endpoint for region cn-hangzhou\n  you need to add --endpoint xxx.aliyuncs.com", err.Error())

	endpointflag.SetAssigned(true)
	err = invooker.Init(ctx, product)
	assert.Nil(t, err)
}
