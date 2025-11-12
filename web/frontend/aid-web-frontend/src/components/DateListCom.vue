<template>
  <q-list>
    <q-item clickable :key="date" v-for="date in dateList">
      <q-item-section>{{ date }}</q-item-section>
    </q-item>
  </q-list>
</template>

<script setup>
defineProps({ dateList: [] })
</script>
