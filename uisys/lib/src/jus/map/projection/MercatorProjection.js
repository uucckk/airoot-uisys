/**
 * 墨卡托经纬度投影
 */
class CommonProjection{
	public Vector2D lonLat2Mercator(Vector2D lonLat)
	{
		Vector2D mercator = new Vector2D();
		double x = lonLat.X * 20037508.34 / 180;
		double y = Math.Log(Math.Tan((90 + lonLat.Y) * Math.PI / 360)) / (Math.PI / 180);
		y = y * 20037508.34 / 180;
		mercator.X = x;
		mercator.Y = y;
		return mercator;
	}
	//墨卡托转经纬度
	public Vector2D Mercator2lonLat(Vector2D mercator)
	{
		Vector2D lonLat = new Vector2D();
		double x = mercator.X / 20037508.34 * 180;
		double y = mercator.Y / 20037508.34 * 180;
		y = 180 / Math.PI * (2 * Math.Atan(Math.Exp(y * Math.PI / 180)) - Math.PI / 2);
		lonLat.X = x;
		lonLat.Y = y;
		return lonLat;
	}
	
}
