(ns plot.core
  (:gen-class)
  (:require [com.hypirion.clj-xchart :as c]
            [clojure.data.csv :as csv]
            [clojure.java.io :as io]
            [clojure.alpha.spec :as s]))

(def regexes-keywords-table
  [[#"[Dd]ate" :date] [#"[Ww]eight" :weight] [#"Points?" :points]
   [#"[Bb]ody[ _][Ff]at" :body-fat-%]
   [#"BMI" :bmi] [#"[Mm]uscle([ _][Mm]ass)?([ _][Ww]eight)?" :muscle]
   [#"Water([ _][Pp]ercent(age)?)?" :water-%]
   [#"Basal[ _][Mm]etabolism" :basal-metabolism]
   [#"Visceral[ _][Ff]at" :visceral-fat]
   [#"Bone[ _][Mm]ass" :bone-mass]
   [#"Sylwetka body( fat)?" :sylwetka-body]
   [#"Sylwetka muscle([ _]ratio)?" :sylwetka-muscle]])

(d)

(defn csvs->vecs-of-rows [csvs-paths]
  (mapv
   (fn csv->vec-of-rows [csv-path]
     (with-open [reader (io/reader csv-path)]
       (let [csv-source (csv/read-csv reader)
             header-kws (loop [done-kws []
                               todo-hs (first csv-source)
                               pasd-rgxs []
                               todo-rgxs regexes-keywords-table]
                          (cond
                            (empty? todo-hs) (if (empty? todo-rgxs) done-kws
                                                 (do (println "Match not found for: " (map second todo-rgxs)) done-kws))
                            (empty? todo-rgxs) (if (empty? pasd-rgxs)
                                                 (do (println "Additional columns: " (pr-str todo-hs)) done-kws)
                                                 (do (println "Unmatched column: " (pr-str (first todo-hs)))
                                                     (recur (conj done-kws (keyword (first todo-hs)))
                                                            (rest todo-hs) [] (into pasd-rgxs todo-rgxs))))
                            (re-matches (first (first todo-rgxs))
                                        (first todo-hs))
                            (recur (conj done-kws (second (first todo-rgxs)))
                                   (rest todo-hs) []
                                   (into pasd-rgxs (rest todo-rgxs)))
                            :else
                            (recur done-kws todo-hs (conj pasd-rgxs (first todo-rgxs)) (rest todo-rgxs))))]
         (mapv zipmap (repeat header-kws) (rest csv-source)))))
   csvs-paths))

(defn -main
  "I don't do a whole lot ... yet."
  [& args]
  (println "Hello, World!"))
