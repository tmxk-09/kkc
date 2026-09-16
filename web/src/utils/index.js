import SparkMD5 from 'spark-md5'

/**
 * 计算文件 MD5 Hash（使用分片读取，避免阻塞主线程）
 * @param {File} file
 * @param {Function} onProgress 进度回调 (0-100)
 * @returns {Promise<string>} MD5 Hash
 */
export function calculateFileMD5(file, onProgress) {
  return new Promise((resolve, reject) => {
    const chunkSize = 2 * 1024 * 1024 // 2MB per chunk for hashing
    const chunks = Math.ceil(file.size / chunkSize)
    let currentChunk = 0
    const spark = new SparkMD5.ArrayBuffer()
    const reader = new FileReader()

    reader.onload = (e) => {
      spark.append(e.target.result)
      currentChunk++
      if (onProgress) {
        onProgress(Math.round((currentChunk / chunks) * 100))
      }
      if (currentChunk < chunks) {
        loadNext()
      } else {
        resolve(spark.end())
      }
    }

    reader.onerror = () => reject(new Error('文件读取失败'))

    function loadNext() {
      const start = currentChunk * chunkSize
      const end = Math.min(start + chunkSize, file.size)
      reader.readAsArrayBuffer(file.slice(start, end))
    }

    loadNext()
  })
}

/**
 * 格式化文件大小
 */
export function formatFileSize(bytes) {
  if (bytes === 0) return '0 B'
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(1024))
  return (bytes / Math.pow(1024, i)).toFixed(2) + ' ' + units[i]
}

/**
 * 格式化速度
 */
export function formatSpeed(bytesPerSecond) {
  return formatFileSize(bytesPerSecond) + '/s'
}
