/***************************************************
 ** @Desc : This file for 交易记录
 ** @Time : 19.12.2 16:34
 ** @Author : Joker
 ** @File : trade_record
 ** @Last Modified by : Joker
 ** @Last Modified time: 19.12.2 16:34
 ** @Software: GoLand
****************************************************/
package controllers

import (
	"dongfeng-pay/jhagent/sys/enum"
	"dongfeng-pay/service/models"
	"strconv"
	"strings"
)

type TradeRecord struct {
	KeepSession
}

// @router /trade/show_ui
func (c *TradeRecord) ShowUI() {
	us := c.GetSession(enum.UserSession)
	u := us.(models.AgentInfo)

	ranMd5 := encrypt.EncodeMd5([]byte(pubMethod.RandomString(46)))
	c.Ctx.SetCookie(enum.UserCookie, ranMd5, enum.CookieExpireTime)
	c.Ctx.SetSecureCookie(ranMd5, enum.UserCookie, ranMd5, enum.CookieExpireTime)
	c.SetSession(enum.UserCookie, ranMd5)

	c.Data["payType"] = enum.GetPayType()
	c.Data["status"] = enum.GetOrderStatus()
	c.Data["userName"] = u.AgentName
	c.TplName = "trade_record.html"
}

// 订单记录查询分页
// @router /trade/list/?:params [get]
func (c *TradeRecord) TradeQueryAndListPage() {
	us := c.GetSession(enum.UserSession)
	u := us.(models.AgentInfo)

	// 分页参数
	page, _ := strconv.Atoi(c.GetString("page"))
	limit, _ := strconv.Atoi(c.GetString("limit"))
	if limit == 0 {
		limit = 15
	}

	// 查询参数
	in := make(map[string]string)
	merchantNo := strings.TrimSpace(c.GetString("MerchantNo"))
	merchantName := strings.TrimSpace(c.GetString("merchantName"))
	start := strings.TrimSpace(c.GetString("start"))
	end := strings.TrimSpace(c.GetString("end"))
	payType := strings.TrimSpace(c.GetString("pay_type"))
	status := strings.TrimSpace(c.GetString("status"))

	in["merchant_order_id"] = merchantNo
	in["merchant_name__icontains"] = merchantName
	in["pay_type_code"] = payType
	in["status"] = status
	in["agent_uid"] = u.AgentUid

	if start != "" {
		in["update_time__gte"] = start
	}
	if end != "" {
		in["update_time_lte"] = end
	}

	// 计算分页数
	count := models.GetOrderProfitLenByMap(in)
	totalPage := count / limit // 计算总页数
	if count%limit != 0 {      // 不满一页的数据按一页计算
		totalPage++
	}

	// 数据获取
	var list []models.OrderProfitInfo
	if page <= totalPage {
		list = models.GetOrderProfitByMap(in, limit, (page-1)*limit)
	}

	// 数据回显
	out := make(map[string]interface{})
	out["limit"] = limit // 分页数据
	out["page"] = page
	out["totalPage"] = totalPage
	out["root"] = list // 显示数据

	c.Data["json"] = out
	c.ServeJSON()
	c.StopRun()
}

// @router /trade/show_complaint_ui
func (c *TradeRecord) ShowComplaintUI() {
	us := c.GetSession(enum.UserSession)
	u := us.(models.AgentInfo)

	ranMd5 := encrypt.EncodeMd5([]byte(pubMethod.RandomString(46)))
	c.Ctx.SetCookie(enum.UserCookie, ranMd5, enum.CookieExpireTime)
	c.Ctx.SetSecureCookie(ranMd5, enum.UserCookie, ranMd5, enum.CookieExpireTime)
	c.SetSession(enum.UserCookie, ranMd5)

	c.Data["payType"] = enum.GetPayType()
	c.Data["status"] = enum.GetComOrderStatus()
	c.Data["userName"] = u.AgentName
	c.TplName = "complaint_record.html"
}

// 投诉列表查询分页
// @router /trade/complaint/?:params [get]
func (c *TradeRecord) ComplaintQueryAndListPage() {
	us := c.GetSession(enum.UserSession)
	u := us.(models.AgentInfo)

	// 分页参数
	page, _ := strconv.Atoi(c.GetString("page"))
	limit, _ := strconv.Atoi(c.GetString("limit"))
	if limit == 0 {
		limit = 15
	}

	// 查询参数
	in := make(map[string]string)
	merchantNo := strings.TrimSpace(c.GetString("MerchantNo"))
	merchantName := strings.TrimSpace(c.GetString("merchantName"))
	start := strings.TrimSpace(c.GetString("start"))
	end := strings.TrimSpace(c.GetString("end"))
	payType := strings.TrimSpace(c.GetString("pay_type"))
	status := strings.TrimSpace(c.GetString("status"))

	in["merchant_order_id"] = merchantNo
	in["merchant_name__icontains"] = merchantName
	in["pay_type_code"] = payType
	if strings.Compare("YES", status) == 0 {
		in["freeze"] = enum.YES
	} else {
		in["refund"] = enum.YES
	}
	in["agent_uid"] = u.AgentUid

	if start != "" {
		in["update_time__gte"] = start
	}
	if end != "" {
		in["update_time__lte"] = end
	}

	// 计算分页数
	count := models.GetOrderLenByMap(in)
	totalPage := count / limit // 计算总页数
	if count%limit != 0 {      // 不满一页的数据按一页计算
		totalPage++
	}

	// 数据获取
	var list []models.OrderInfo
	if page <= totalPage {
		list = models.GetOrderByMap(in, limit, (page-1)*limit)
	}

	// 数据回显
	out := make(map[string]interface{})
	out["limit"] = limit // 分页数据
	out["page"] = page
	out["totalPage"] = totalPage
	out["root"] = list // 显示数据

	c.Data["json"] = out
	c.ServeJSON()
	c.StopRun()
}
