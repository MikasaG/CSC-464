package matrix;

import java.util.ArrayList;
import java.util.List;
import java.util.Random;
import java.util.concurrent.Callable;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.concurrent.Future;

public class Matrix_Linear {
	private int[][] value;

    private Matrix_Linear(int[][] value) {
        this.value = value;
    }

    Matrix_Linear(int rows, int cols) {
        Random r = new Random(0);
        value = new int[rows][cols];
        for (int i = 0; i < rows; i++) {
            for (int j = 0; j < cols; j++) {
                this.value[i][j] = r.nextInt();
            }
        }
    }

    int[][] getValue() {
        return value;
    }

    private int rows() {
        return value.length;
    }

    private int cols() {
        return value[0].length;
    }

    Matrix_Linear multiply(Matrix_Linear other, int parallelism) {

        int[][] result = new int[rows()][cols()];
        ExecutorService executor = Executors.newFixedThreadPool(parallelism);
        List<Future<int[][]>> callables = new ArrayList<>();

        int parts = Math.max(rows() / parallelism, 1);

        for (int i = 0; i < rows(); i += parts) {
            Callable<int[][]> worker = new PartMultiplier(value, other.getValue(), i, i + parts);
            callables.add(executor.submit(worker));
        }

        // retrieve the result from parts
        int start = 0;
        for (Future<int[][]> future : callables) {
            try {
                System.arraycopy(future.get(), start, result, start, parts);
            } catch (InterruptedException | ExecutionException e) {
                e.printStackTrace();
            }
            start += parts;
        }
        executor.shutdown();

        return new Matrix_Linear(result);
    }

    private class PartMultiplier implements Callable<int[][]> {

        int[][] a, b, c;
        int start, end;

        PartMultiplier(int[][] a, int[][] b, int s, int e) {
            this.a = a;
            this.b = b;
            this.c = new int[a.length][b[0].length];
            start = s;
            end = e;
        }

        @Override
        public int[][] call() {
            System.err.println(Thread.currentThread().getName());
            for (int i = start; i < end; i++) {
                for (int k = 0; k < b.length; k++) {
                    for (int j = 0; j < b[0].length; j++) {
                        c[i][j] += a[i][k] * b[k][j];
                    }
                }
            }
            return c;
        }
    }

    @Override
    public String toString() {
        StringBuilder sb = new StringBuilder(value.length);
        for (int[] row : value) {
            int i = 0;
            for (int number : row) {
                if (i != 0) {
                    sb.append("\t");
                } else {
                    i++;
                }
                sb.append(number);
            }
            sb.append("\n\r");
        }
        return sb.toString();
    }
    
    public static void main(String[] args) {

        // For test data
     /*   Matrix_Linear m1 = new Matrix_Linear(10,10);
        Matrix_Linear m2 = new Matrix_Linear(10,10);
        int threads = Runtime.getRuntime().availableProcessors();

        System.out.println(m1);
        System.out.println(m2);

        long start = System.nanoTime();
        Matrix_Linear m3 = m1.multiply(m2, threads);
        long end = System.nanoTime();
        double seconds = (double)(end - start) / 1000000000.0;

        System.out.println(m3);
        System.out.println("Time: " + seconds + " s. Number of treads: " + threads);*/

        // For big matrices with random data
        timedTest(1000, 1);

        timedTest(2000, 1);

        timedTest(4000, 1);
    }

    private static void timedTest(int dimension, int parallelism) {

        Matrix_Linear mr1 = new Matrix_Linear(dimension, dimension);
        Matrix_Linear mr2 = new Matrix_Linear(dimension, dimension);

        long start = System.nanoTime();
        long beforeUsedMem=Runtime.getRuntime().totalMemory()-Runtime.getRuntime().freeMemory();
        int dim = mr1.multiply(mr2, parallelism).getValue().length;
		long afterUsedMem=Runtime.getRuntime().totalMemory()-Runtime.getRuntime().freeMemory();
        long end = System.nanoTime();
        double seconds = (double)(end - start) / 1000000000.0;
        System.out.println("Dim: " + dim + ", time: " + seconds + " s. Memory used: " + (afterUsedMem-beforeUsedMem)/(1024.0 * 1024.0)+" MB. Parallelism: " + parallelism);
    }
}

