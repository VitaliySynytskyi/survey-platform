<template>
  <div class="not-found-page">
    <div class="container text-center">
      <h1>404</h1>
      <h2>Page Not Found</h2>
      <p>The page you are looking for does not exist or has been moved.</p>
      <router-link to="/" class="btn">Go to Home</router-link>
    </div>
  </div>
</template>

<script setup lang="ts">
// No additional script needed
</script>

<style scoped>
.not-found-page {
  padding: 4rem 0;
  text-align: center;
}

.not-found-page h1 {
  font-size: 6rem;
  margin-bottom: 1rem;
  color: var(--primary-color);
}

.not-found-page h2 {
  margin-bottom: 1.5rem;
  color: var(--text-color);
}

.not-found-page p {
  margin-bottom: 2rem;
  color: var(--dark-gray);
  font-size: 1.1rem;
}

.not-found-page .btn {
  display: inline-block;
  margin-top: 1rem;
}
</style> 