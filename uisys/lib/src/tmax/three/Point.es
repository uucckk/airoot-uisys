class Point{
	function init(){
		sprite = new ? THREE.TextureLoader().load( 'img/Frame/jingli_P.png' );
			var sGeometry = new ?THREE.BufferGeometry();
			var sPositions = [];
			var num = 1;
			var r = 10;
			var radian = Math.PI / (18 * num);
			var sMaterial = new ?THREE.PointsMaterial({
						size: 0.5,
						transparent: true,
						color: "#ffffff",
						alphaTest: 0.5,
						size:1,
						map:sprite
						});
			for (var i = 0; i < 36 * num; i+=3) {
				sPositions[i] = r * Math.sin(radian * i);
				sPositions[i+1] = 0;
				sPositions[i+2] = r * Math.cos(radian * i);                
			}            
			sGeometry.addAttribute('position', new ?THREE.Float32BufferAttribute(sPositions, 3));
			sGeometry.computeBoundingSphere();
			var orbit = new ?THREE.Points(sGeometry, sMaterial);
			scene.add(orbit);
	}
}