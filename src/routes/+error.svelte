<script>
	import { pb } from '$lib/pb';
	import PageWrapper from '$lib/shell/PageWrapper.svelte';
	import Oops from '$lib/assets/oops.webp?w=400;300&enhanced';
</script>

<PageWrapper title="Cubby" {pb}>
	<div class="grid h-full content-center justify-items-center gap-4">
		<enhanced:img
			src={Oops}
			alt="Error!"
			sizes="(min-width:1920px) 400px, (min-width:100px) 300px"
		/>
		<h2 class="text-xl font-medium">Ugh, something's broken.</h2>
	</div>
</PageWrapper>
