package infra

import "github.com/oceanho/gw"

var PermissionCreateTestDataDecorator = gw.NewPermissionDecorator(PermissionCreateTestData)

var DecoratorList = gw.NewPermAllDecorator("Tester")
