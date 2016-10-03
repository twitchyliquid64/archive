class RollingAverage {
  float curVal;
  float[] vals = new float[20];
  int count;
  float ad(float val)
  {
    vals[count%vals.length] = val;
    count++;
    return val();
  }
  
  float val()
  {
    float ret = 0;
    for(float v: vals)
      ret += v;
    return ret/(vals.length+1);
  }
}
