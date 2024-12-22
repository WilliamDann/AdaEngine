namespace Chess;

public class Point
{
    public int x;
    public int y;

    public Point(int x, int y)
    {
        this.x = x;
        this.y = y;
    }

    public Point(string notation)
    {
        this.x = (char)(notation[0] - 97);
        this.y = notation[1];
    }

    public override string ToString()
    {
        return $"{(char)(97 + (7 - this.x))}{this.y + 1}";
    }
}