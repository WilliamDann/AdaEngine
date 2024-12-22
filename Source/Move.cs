namespace Chess;

public class Move
{
    public Point from;
    public Point to;

    public Move(Point from, Point to)
    {
        this.from = from;
        this.to   = to;
    }

    // uci notation of the move
    public string UCI()
    {
        return $"{this.from}{this.to}";
    }

    // TODO San notaiton, requires board info for captures?

    public override string ToString()
    {
        return this.UCI();
    }
}