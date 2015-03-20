import System.IO

-- sum' :: (Num a) => [a] -> a
sum' xs = foldl (\acc x -> acc + x) 0 xs

main = do
    putStrLn $ show $ sum' [3, 5, 2, 1]
