(def project 'archiet/mmih.plot)
(def version "0.1.0-SNAPSHOT")

(set-env! :resource-paths #{"resources" "src"}
          :source-paths   #{"test"}
          :dependencies   '[[org.clojure/clojure "1.9.0-alpha17"]
                            [adzerk/boot-test "RELEASE" :scope "test"]
                            [com.hypirion/clj-xchart "0.2.0"]
                            ;;[org.clojure/tools.cli "0.3.5"]
                            [org.clojure/data.csv "0.1.4"]])

(task-options!
 aot {:namespace   #{'archiet/mmih.plot.core}}
 pom {:project     project
      :version     version
      :description "plotting data from mmih"
      :url         "http://github.com/ArchieT/miscale-manual-input-helper"
      :scm         {:url "https://github.com/ArchieT/miscale-manual-input-helper"}}
 jar {:main        'archiet/mmih.plot.core
      :file        (str "plot-mmih-" version "-standalone.jar")})

(deftask build
  "Build the project locally as a JAR."
  [d dir PATH #{str} "the set of directories to write to (target)."]
  (let [dir (if (seq dir) dir #{"target"})]
    (comp (aot) (pom) (uber) (jar) (target :dir dir))))

(deftask run
  "Run the project."
  [a args ARG [str] "the arguments for the application."]
  (require '[plot.core :as app])
  (apply (resolve 'app/-main) args))

(require '[adzerk.boot-test :refer [test]])
