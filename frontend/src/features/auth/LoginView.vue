<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { useUserStore } from '@/stores/user'
import { resetDynamicRoutes } from '@/app/router'

const router = useRouter()
const route = useRoute()
const user = useUserStore()
const formRef = ref<FormInstance>()
const loading = ref(false)

const form = reactive({
  account_number: 'admin',
  password: 'admin',
})

const rules: FormRules = {
  account_number: [{ required: true, message: '请输入账号', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
}

async function onSubmit() {
  await formRef.value?.validate()
  loading.value = true
  try {
    resetDynamicRoutes()
    await user.login(form.account_number, form.password)
    ElMessage.success('登录成功')
    const redirect = (route.query.redirect as string) || '/'
    await router.replace(redirect)
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="login-page">
    <div class="panel">
      <div class="panel-head">
        <h1>RBAC 权限管理</h1>
        <p>基于角色的访问控制后台</p>
      </div>
      <el-form ref="formRef" :model="form" :rules="rules" size="large" @keyup.enter="onSubmit">
        <el-form-item prop="account_number">
          <el-input v-model="form.account_number" placeholder="账号" clearable />
        </el-form-item>
        <el-form-item prop="password">
          <el-input v-model="form.password" type="password" placeholder="密码" show-password />
        </el-form-item>
        <el-button type="primary" class="submit" :loading="loading" @click="onSubmit">
          登 录
        </el-button>
      </el-form>
    </div>
  </div>
</template>

<style scoped>
.login-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background:
    radial-gradient(circle at 20% 20%, rgba(64, 158, 255, 0.18), transparent 40%),
    radial-gradient(circle at 80% 0%, rgba(103, 194, 58, 0.12), transparent 35%),
    linear-gradient(160deg, #1f2d3d 0%, #2c3e50 45%, #1a2332 100%);
}
.panel {
  width: 400px;
  padding: 36px 32px 28px;
  border-radius: 12px;
  background: rgba(255, 255, 255, 0.96);
  box-shadow: 0 16px 48px rgba(0, 0, 0, 0.25);
}
.panel-head {
  text-align: center;
  margin-bottom: 28px;
}
.panel-head h1 {
  margin: 0;
  font-size: 24px;
  color: #1f2d3d;
  font-weight: 600;
}
.panel-head p {
  margin: 8px 0 0;
  color: #909399;
  font-size: 13px;
}
.submit {
  width: 100%;
  margin-top: 8px;
}
</style>
