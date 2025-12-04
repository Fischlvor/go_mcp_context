/**
 * 生成稳定的浏览器指纹（设备ID）
 * 基于多个浏览器特征生成，即使清除localStorage也能保持一致
 */

// 简单的字符串哈希函数
function hashString(str: string): string {
  let hash = 0
  for (let i = 0; i < str.length; i++) {
    const char = str.charCodeAt(i)
    hash = ((hash << 5) - hash) + char
    hash = hash & hash // Convert to 32bit integer
  }
  return Math.abs(hash).toString(36)
}

// 收集浏览器特征
function collectBrowserFingerprint(): string {
  const features: string[] = []
  
  // 1. User-Agent（稳定）
  features.push(navigator.userAgent || '')
  
  // 2. 屏幕分辨率（稳定）
  features.push(`${screen.width}x${screen.height}`)
  features.push(`${screen.availWidth}x${screen.availHeight}`)
  
  // 3. 颜色深度（稳定）
  features.push((screen.colorDepth || '').toString())
  
  // 4. 时区（稳定）
  features.push(Intl.DateTimeFormat().resolvedOptions().timeZone || '')
  
  // 5. 语言（稳定）
  features.push(navigator.language || '')
  features.push((navigator.languages || []).join(','))
  
  // 6. 硬件并发数（稳定）
  features.push((navigator.hardwareConcurrency || '').toString())
  
  // 7. 平台信息（稳定）
  features.push(navigator.platform || '')
  
  // 10. Canvas指纹（简化版，不实际绘制）
  if (typeof document !== 'undefined') {
    try {
      const canvas = document.createElement('canvas')
      const ctx = canvas.getContext('2d')
      if (ctx) {
        ctx.textBaseline = 'top'
        ctx.font = '14px Arial'
        ctx.fillText('Browser fingerprint', 2, 2)
        // 不实际获取像素数据，只使用canvas特性
        features.push(canvas.width.toString() + canvas.height.toString())
      }
    } catch (e) {
      // Canvas不可用
    }
  }
  
  // 11. WebGL信息（如果可用）
  if (typeof document !== 'undefined') {
    try {
      const canvas = document.createElement('canvas')
      const gl = canvas.getContext('webgl') || canvas.getContext('experimental-webgl') as WebGLRenderingContext
      if (gl) {
        const debugInfo = gl.getExtension('WEBGL_debug_renderer_info')
        if (debugInfo) {
          features.push(gl.getParameter(debugInfo.UNMASKED_VENDOR_WEBGL) || '')
          features.push(gl.getParameter(debugInfo.UNMASKED_RENDERER_WEBGL) || '')
        }
      }
    } catch (e) {
      // WebGL不可用
    }
  }
  
  // 12. 字体检测（简化版，检测常见字体）
  // 注意：需要在浏览器环境中执行，不能在SSR时执行
  if (typeof document !== 'undefined' && document.body) {
    try {
      const fonts = ['Arial', 'Verdana', 'Times New Roman', 'Courier New', 'Georgia']
      const fontAvailable: string[] = []
      const testString = 'mmmmmmmmmmlli'
      const testSize = '72px'
      
      const span = document.createElement('span')
      span.style.fontSize = testSize
      span.innerHTML = testString
      span.style.position = 'absolute'
      span.style.left = '-9999px'
      span.style.top = '-9999px'
      document.body.appendChild(span)
      
      const defaultWidth = span.offsetWidth
      const defaultHeight = span.offsetHeight
      
      fonts.forEach(font => {
        span.style.fontFamily = font
        const width = span.offsetWidth
        const height = span.offsetHeight
        if (width !== defaultWidth || height !== defaultHeight) {
          fontAvailable.push(font)
        }
      })
      
      document.body.removeChild(span)
      features.push(fontAvailable.join(','))
    } catch (e) {
      // 字体检测失败，跳过
    }
  }
  
  return features.join('|')
}

/**
 * 生成或获取设备ID
 * 优先从localStorage读取，如果不存在则生成浏览器指纹
 * 即使清除localStorage，只要浏览器环境不变，指纹也会保持一致
 */
export function getDeviceId(): string {
  // 1. 优先从localStorage读取
  let deviceId = localStorage.getItem('device_id')
  
  if (deviceId) {
    return deviceId
  }
  
  // 2. 如果不存在，生成浏览器指纹
  try {
    const fingerprint = collectBrowserFingerprint()
    // 生成哈希值作为设备ID
    deviceId = 'fp_' + hashString(fingerprint)
    
    // 存储到localStorage（用于下次快速读取）
    localStorage.setItem('device_id', deviceId)
    
    console.log('✓ 设备ID已生成:', deviceId)
    return deviceId
  } catch (e) {
    // 如果指纹生成失败，使用随机ID作为后备方案
    console.warn('生成浏览器指纹失败，使用随机ID:', e)
    deviceId = 'web_' + Math.random().toString(36).substring(2, 15) + Date.now().toString(36)
    localStorage.setItem('device_id', deviceId)
    return deviceId
  }
}
