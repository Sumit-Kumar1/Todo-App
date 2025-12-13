<div class="block404">
	<div class="waves"></div>
	<div class="obj">
		<img src="https://imgur.com/w0Yb4MX.png" alt="" />
		<a class="px-10 text-2xl text-black" href="/">home</a>
	</div>
	<div class="t404"></div>

	<svg xmlns="http://www.w3.org/2000/svg" version="1.1">
		<defs>
			<filter id="glitch">
				<feTurbulence
					type="fractalNoise"
					baseFrequency="0.01 0.03"
					numOctaves="1"
					result="warp"
					id="turb"
				/>
				<feColorMatrix in="warp" result="huedturb" type="hueRotate" values="90">
					<animate
						attributeType="XML"
						attributeName="values"
						values="0;180;360"
						dur="3s"
						repeatCount="indefinite"
					/>
				</feColorMatrix>
				<feDisplacementMap
					xChannelSelector="R"
					yChannelSelector="G"
					scale="50"
					in="SourceGraphic"
					in2="huedturb"
				/>
			</filter>
		</defs>
	</svg>
</div>

<style>
	.block404 {
		position: relative;
		background-image: url(https://imgur.com/eHVyWTM.jpg);
		width: 100%;
		height: 1024px;
		overflow: hidden;
		cursor: pointer;
	}

	.t404 {
		position: absolute;
		width: 364px;
		height: 146px;
		top: 50%;
		left: 50%;
		transform: translate(-50%, -50%);
		background-image: url(https://imgur.com/KPZo9YX.png);
	}

	.obj {
		width: 204px;
		height: 209px;
		position: absolute;
		top: 50%;
		left: 50%;
		transform: translate(-50%, -50%);
		animation: animation-404 6s infinite ease-in-out;
	}

	@keyframes animation-404 {
		0%,
		100% {
			transform: translate(-50%, -50%) rotate(0);
		}
		50% {
			transform: translate(-53%, -42%) rotate(-5deg);
		}
	}

	.waves {
		position: absolute;
		width: 100%;
		height: 100%;
		background-image: url(https://imgur.com/eHVyWTM.jpg);
		filter: url('#glitch');
	}
</style>
