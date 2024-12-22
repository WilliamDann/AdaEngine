namespace Engine;
using Chess;

public class Program
{
    public static void Main()
    {
        Chessboard b     = new Chessboard(FEN.INITAL);
        MoveGenerator mv = new MoveGenerator(b);

        Console.WriteLine(b);

        List<Move> moves = mv.Moves();
        foreach (Move move in moves)
            Console.WriteLine(move);
        Console.WriteLine(moves.Count);

    }
}