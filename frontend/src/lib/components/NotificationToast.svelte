<script lang="ts">
	import { notifications, type ToastNotification } from '$lib/stores/notifications';

	function iconAndColor(level: ToastNotification['level']) {
		switch (level) {
			case 'success':
				return { bg: 'alert-success', icon: 'M5 13l4 4L19 7' };
			case 'warning':
				return { bg: 'alert-warning', icon: 'M12 9v2m0 4h.01M12 2l10 18H2L12 2z' };
			case 'error':
				return { bg: 'alert-error', icon: 'M6 18L18 6M6 6l12 12' };
			default:
				return {
					bg: 'alert-info',
					icon: 'M13 16h-1v-4h-1m1-4h.01M12 2a10 10 0 100 20 10 10 0 000-20z'
				};
		}
	}
</script>

<div class="pointer-events-none fixed inset-0 z-50">
	<div class="absolute top-4 right-4 flex flex-col items-end gap-2">
		{#each $notifications as n (n.id)}
			{@const style = iconAndColor(n.level)}
			<div
				class={`pointer-events-auto alert w-80 max-w-[90vw] shadow-lg ${style.bg}`}
				role="status"
				aria-live="polite"
			>
				<svg class="h-5 w-5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d={style.icon} />
				</svg>
				<span class="text-sm">{n.message}</span>
			</div>
		{/each}
	</div>
</div>
