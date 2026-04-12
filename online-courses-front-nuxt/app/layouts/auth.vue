<template>
  <div class="auth-shell">
    <div class="auth-shell__backdrop" />
    <main class="auth-shell__main">
      <slot />
    </main>
  </div>
</template>

<style scoped>
.auth-shell {
  min-height: 100vh;
  display: grid;
  place-items: center;
  padding: 1rem;
  position: relative;
  overflow: hidden;
  background: linear-gradient(180deg, #e2e8f0 0%, #f8fafc 100%);
}

.auth-shell__backdrop {
  position: absolute;
  inset: 0;
  background:
    radial-gradient(circle at 15% 20%, rgba(14, 165, 233, 0.22), transparent 20rem),
    radial-gradient(circle at 85% 15%, rgba(249, 115, 22, 0.2), transparent 18rem),
    radial-gradient(circle at 50% 90%, rgba(59, 130, 246, 0.12), transparent 22rem);
}

.auth-shell__main {
  position: relative;
  z-index: 1;
  width: min(980px, 100%);
}
</style>
