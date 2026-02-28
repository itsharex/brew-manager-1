import { reactive, ref, computed } from 'vue'
import { GetBrewData, StartService, StopService, GetAppIcon  } from '../../wailsjs/go/main/App'

export function useBrew() {
  const data = reactive({ formulae: [], casks: [], loading: false })
  const searchQuery = ref('')
  const processingMap = reactive(new Map())
  const toast = reactive({ show: false, msg: '', type: 'success' })

  // 统一的提示函数
  function showToast(msg, type = 'success') {
    toast.msg = msg
    toast.type = type
    toast.show = true
    setTimeout(() => { toast.show = false }, 3000)
  }

  // 搜索过滤逻辑
  const filteredFormulae = computed(() => 
    data.formulae.filter(item => item.name.toLowerCase().includes(searchQuery.value.toLowerCase()))
  )

  const filteredCasks = computed(() => 
    data.casks.filter(item => item.name.toLowerCase().includes(searchQuery.value.toLowerCase()))
  )

  // 数据刷新逻辑
  async function updateList() {
    try {
      const res = await GetBrewData()
      // 遍历 Casks 为每个应用请求图标（或者在后端直接一次性返回）
      const caskWithIcons = await Promise.all(res.casks.map(async (item) => {
        const icon = await GetAppIcon(item.name)
        return { ...item, iconBase64: icon }
      }))
      data.casks = caskWithIcons
      data.formulae = res.formulae
      // data.casks = res.casks
    } catch (err) {
      console.error("刷新失败:", err)
    }
  }

  // 服务启动/停止逻辑
  async function handleService(item) {
    processingMap.set(item.name, true)
    try {
      let result = item.status === 'started' ? await StopService(item.name) : await StartService(item.name)
      if (result.success) {
        showToast(result.message, 'success')
        await updateList()
      } else {
        showToast(result.message, 'error')
      }
    } catch (err) {
      showToast("系统错误: " + err, 'error')
    } finally {
      processingMap.delete(item.name)
    }
  }
  return {
    data, searchQuery, processingMap, toast,
    filteredFormulae, filteredCasks,
    updateList, handleService,
    showToast
  }
}