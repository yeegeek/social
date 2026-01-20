// Package middleware 提供安全响应头中间件
package middleware

import (
	"github.com/gin-gonic/gin"
)

// SecurityHeaders 添加安全响应头中间件
// 设置各种安全相关的 HTTP 响应头，防止常见的 Web 安全漏洞
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 防止 MIME 类型嗅探
		// 告诉浏览器不要尝试猜测内容类型，严格按照 Content-Type 处理
		c.Header("X-Content-Type-Options", "nosniff")

		// 防止点击劫持攻击
		// 禁止页面被嵌入到 iframe 中
		c.Header("X-Frame-Options", "DENY")

		// 启用浏览器的 XSS 过滤器
		// 当检测到 XSS 攻击时，阻止页面加载
		c.Header("X-XSS-Protection", "1; mode=block")

		// 强制使用 HTTPS
		// 告诉浏览器在接下来的一年内，只能通过 HTTPS 访问此站点
		// includeSubDomains: 该规则也适用于所有子域名
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

		// 内容安全策略 (CSP)
		// 限制资源的加载来源，防止 XSS 攻击
		// default-src 'self': 默认只允许加载同源资源
		// script-src 'self' 'unsafe-inline': 允许同源脚本和内联脚本
		// style-src 'self' 'unsafe-inline': 允许同源样式和内联样式
		// img-src 'self' data: https:: 允许同源图片、data URI 和 HTTPS 图片
		c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self' data:; connect-src 'self'")

		// 引用者策略
		// 控制 Referer 头的发送行为
		// no-referrer-when-downgrade: 从 HTTPS 到 HTTP 时不发送 Referer
		c.Header("Referrer-Policy", "no-referrer-when-downgrade")

		// 权限策略
		// 控制浏览器功能的使用权限
		// 禁用地理位置、麦克风、摄像头等敏感功能
		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

		c.Next()
	}
}

// SecurityHeadersCustom 创建自定义安全响应头中间件
// 允许自定义 CSP 策略和其他安全头
func SecurityHeadersCustom(csp string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Header("Referrer-Policy", "no-referrer-when-downgrade")
		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

		// 使用自定义 CSP 策略
		if csp != "" {
			c.Header("Content-Security-Policy", csp)
		}

		c.Next()
	}
}
