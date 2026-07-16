<script lang="ts">
	import '../app.css';
	import favicon from '$lib/assets/favicon.svg';
	import NotificationToast from '$lib/components/NotificationToast.svelte';
	import { theme } from '$lib/stores/theme';
	import { onMount } from 'svelte';

	let { children } = $props();

	onMount(() => {
		document.documentElement.setAttribute('data-theme', $theme);
	});
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

{@render children?.()}

<NotificationToast />
