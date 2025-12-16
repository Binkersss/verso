(function() {
	let animationRunning = false;

	function initParticles() {
		const c = document.getElementById('particles-canvas');
		if (!c || animationRunning) {
			return;
		}

		animationRunning = true;
		const ctx = c.getContext('2d');

		c.width = c.offsetWidth || 800;
		c.height = 400;

		console.log('Particles initialized!');

		const particles = [];
		const particleCount = 100;
		const connectionDist = 120;
		const mouse = { x: null, y: null, radius: 150 };

		class Particle {
			constructor() {
				this.startX = Math.random() * c.width;
				this.startY = Math.random() * c.height;
				this.x = this.startX;
				this.y = this.startY;
				this.vx = (Math.random() - 0.5) * 1.5;
				this.vy = (Math.random() - 0.5) * 1.5;
				this.size = Math.random() * 3 + 2;
			}

			update() {
				if (this.x > c.width || this.x < 0) this.vx *= -1;
				if (this.y > c.height || this.y < 0) this.vy *= -1;

				const dx = mouse.x - this.x;
				const dy = mouse.y - this.y;
				const dist = Math.sqrt(dx * dx + dy * dy);

				if (dist < mouse.radius && mouse.x !== null) {
					const force = (mouse.radius - dist) / mouse.radius;
					const angle = Math.atan2(dy, dx);
					this.vx -= Math.cos(angle) * force * 0.3;
					this.vy -= Math.sin(angle) * force * 0.3;
				}

				// Gentle drift back to starting position
				const dxStart = this.startX - this.x;
				const dyStart = this.startY - this.y;
				this.vx += dxStart * 0.002;
				this.vy += dyStart * 0.002;

				this.x += this.vx;
				this.y += this.vy;

				this.vx *= 0.98;
				this.vy *= 0.98;
			}

			draw() {
				ctx.fillStyle = '#63b3ed';
				ctx.beginPath();
				ctx.arc(this.x, this.y, this.size, 0, Math.PI * 2);
				ctx.fill();
			}
		}

		for (let i = 0; i < particleCount; i++) {
			particles.push(new Particle());
		}

		function connect() {
			for (let i = 0; i < particles.length; i++) {
				for (let j = i + 1; j < particles.length; j++) {
					const dx = particles[i].x - particles[j].x;
					const dy = particles[i].y - particles[j].y;
					const dist = Math.sqrt(dx * dx + dy * dy);

					if (dist < connectionDist) {
						const opacity = (1 - dist / connectionDist) * 0.5;
						ctx.strokeStyle = `rgba(99, 179, 237, ${opacity})`;
						ctx.lineWidth = 1;
						ctx.beginPath();
						ctx.moveTo(particles[i].x, particles[i].y);
						ctx.lineTo(particles[j].x, particles[j].y);
						ctx.stroke();
					}
				}
			}
		}

		function animate() {
			if (!document.getElementById('particles-canvas')) {
				animationRunning = false;
				return;
			}

			ctx.fillStyle = '#0a0e27';
			ctx.fillRect(0, 0, c.width, c.height);

			particles.forEach(p => {
				p.update();
				p.draw();
			});

			connect();
			requestAnimationFrame(animate);
		}

		c.addEventListener('mousemove', e => {
			const rect = c.getBoundingClientRect();
			mouse.x = e.clientX - rect.left;
			mouse.y = e.clientY - rect.top;
		});

		c.addEventListener('mouseleave', () => {
			mouse.x = null;
			mouse.y = null;
		});

		animate();
	}

	// Check for canvas periodically
	setInterval(initParticles, 500);
})();
