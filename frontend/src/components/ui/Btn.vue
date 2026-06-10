<template>
  <button
    class="btn"
    :class="[variant ? `btn-${variant}` : '', { 'btn-sm': sm, 'btn-icon': icon }]"
    :disabled="disabled || loading"
    :type="type"
  >
    <span v-if="loading" class="btn-spinner" />
    <slot v-else />
  </button>
</template>

<script lang="ts" setup>
withDefaults(defineProps<{
  variant?: 'primary' | 'ghost' | 'subtle' | ''
  sm?: boolean
  icon?: boolean
  loading?: boolean
  disabled?: boolean
  type?: 'button' | 'submit'
}>(), {
  variant: 'ghost',
  type: 'button',
})
</script>

<style scoped>
.btn-spinner {
  width: 16px; height: 16px;
  border: 2px solid rgba(255, 255, 255, .4);
  border-top-color: #fff;
  border-radius: 50%;
  animation: spin .7s linear infinite;
}
.btn-ghost .btn-spinner, .btn-subtle .btn-spinner {
  border-color: var(--line-2);
  border-top-color: var(--text-2);
}
</style>
