<script lang="ts">
	let { children } = $props();

	const removeTopMarginForContentCenter = true;
</script>

<div class={[removeTopMarginForContentCenter && 'lg:-pt-12 lg:-mt-14']}>
	{@render children()}
</div>
