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
    <aside class="login-visual" aria-hidden="true">
      <!-- 预留大图位：后续可换成 <img> 或背景图 -->
      <div class="visual-stage">
        <p class="brand">RBAC</p>
        <p class="brand-sub">权限管理后台</p>
      </div>
    </aside>

    <section class="login-side">
      <div class="login-panel">
        <header class="panel-head">
          <h1>登录</h1>
          <p>使用账号密码进入系统</p>
        </header>
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
    </section>
  </div>
</template>

<style scoped>
.login-page {
  --login-ink: #0f1c24;
  --login-accent: #2a6f7a;
  min-height: 100vh;
  display: grid;
  grid-template-columns: minmax(0, 1.35fr) minmax(360px, 460px);
  background: #e8eef0;
}

.login-visual {
  position: relative;
  min-height: 100vh;
  overflow: hidden;
  background:
    linear-gradient(145deg, rgba(15, 28, 36, 0.55) 0%, rgba(42, 111, 122, 0.35) 55%, rgba(15, 28, 36, 0.7) 100%),
    linear-gradient(160deg, #1a3340 0%, #2a6f7a 48%, #0f1c24 100%);
}

.visual-stage {
  position: absolute;
  inset: 0;
  display: flex;
  flex-direction: column;
  justify-content: flex-end;
  padding: 48px 56px 64px;
  /* 后续大图示例：
     background: url('/login-hero.jpg') center / cover no-repeat;
  */
}

.brand {
  margin: 0;
  font-family: 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei', sans-serif;
  font-size: clamp(48px, 7vw, 88px);
  font-weight: 700;
  letter-spacing: 0.08em;
  line-height: 1;
  color: #f4f8f9;
  text-shadow: 0 8px 32px rgba(0, 0, 0, 0.25);
  animation: brand-in 700ms ease-out both;
}

.brand-sub {
  margin: 14px 0 0;
  font-size: 16px;
  letter-spacing: 0.28em;
  color: rgba(244, 248, 249, 0.78);
  animation: brand-in 700ms ease-out 120ms both;
}

.login-side {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 40px 32px;
  background: #f7fafb;
  box-shadow: -12px 0 40px rgba(15, 28, 36, 0.08);
  animation: side-in 560ms ease-out both;
}

.login-panel {
  width: 100%;
  max-width: 360px;
}

.panel-head {
  margin-bottom: 28px;
}

.panel-head h1 {
  margin: 0;
  font-size: 28px;
  font-weight: 600;
  color: var(--login-ink);
  letter-spacing: 0.04em;
}

.panel-head p {
  margin: 10px 0 0;
  color: #6b7c86;
  font-size: 14px;
}

.submit {
  width: 100%;
  margin-top: 8px;
  --el-button-bg-color: var(--login-accent);
  --el-button-border-color: var(--login-accent);
  --el-button-hover-bg-color: #245f68;
  --el-button-hover-border-color: #245f68;
}

@keyframes brand-in {
  from {
    opacity: 0;
    transform: translateY(18px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes side-in {
  from {
    opacity: 0;
    transform: translateX(24px);
  }
  to {
    opacity: 1;
    transform: translateX(0);
  }
}

@media (max-width: 900px) {
  .login-page {
    grid-template-columns: 1fr;
    grid-template-rows: 220px 1fr;
  }

  .login-visual {
    min-height: 220px;
  }

  .visual-stage {
    justify-content: center;
    padding: 32px 28px;
  }

  .brand {
    font-size: 40px;
  }

  .login-side {
    box-shadow: none;
    align-items: flex-start;
    padding-top: 36px;
  }
}
</style>
